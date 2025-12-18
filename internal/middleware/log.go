package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
)

// RequestLog 请求日志结构
type RequestLog struct {
	Timestamp    string `json:"timestamp"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Query        string `json:"query,omitempty"`
	ClientIP     string `json:"client_ip"`
	UserAgent    string `json:"user_agent,omitempty"`
	StatusCode   int    `json:"status_code"`
	Latency      string `json:"latency"`
	ResponseSize int    `json:"response_size"`
	Error        string `json:"error,omitempty"`
	RequestBody  string `json:"request_body,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
}

// GinLogger 结构化日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 读取请求体（对于POST/PUT请求）
		var requestBody string
		if shouldLogRequestBody(c.Request.Method, c.Request.Header.Get("Content-Type")) {
			requestBody = readRequestBody(c)
		}

		// 处理请求
		c.Next()

		// 记录日志
		logEntry := RequestLog{
			Timestamp:    start.Format(time.DateTime),
			Method:       c.Request.Method,
			Path:         path,
			Query:        raw,
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			StatusCode:   c.Writer.Status(),
			Latency:      time.Since(start).String(),
			ResponseSize: c.Writer.Size(),
			RequestBody:  requestBody,
			RequestID:    requestID,
		}

		// 检查是否有错误
		if len(c.Errors) > 0 {
			logEntry.Error = c.Errors.String()
		}

		// 根据状态码选择日志级别
		logJSON, _ := json.Marshal(logEntry)
		logStr := string(logJSON)

		switch {
		case logEntry.StatusCode >= 500:
			slog.Error("HTTP Request", "details", logStr)
		case logEntry.StatusCode >= 400:
			slog.Warn("HTTP Request", "details", logStr)
		case time.Since(start) > 2*time.Second: // 慢请求
			slog.Warn("Slow HTTP Request", "details", logStr)
		default:
			slog.Info("HTTP Request", "details", logStr)
		}
	}
}

// GinRecovery 恢复中间件增强版
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否为断开的连接
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 结构化故障日志
				panicLog := map[string]interface{}{
					"timestamp":  time.Now().Format(time.DateTime),
					"error":      fmt.Sprintf("%v", err),
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
					"client_ip":  c.ClientIP(),
					"user_agent": c.Request.UserAgent(),
					"request":    string(httpRequest),
				}

				if requestID, exists := c.Get("request_id"); exists {
					panicLog["request_id"] = requestID
				}

				if stack {
					panicLog["stack"] = string(debug.Stack())
				}

				logJSON, _ := json.Marshal(panicLog)

				if brokenPipe {
					slog.Error("Broken pipe panic", "details", string(logJSON))
					c.Error(err.(error))
					c.Abort()
					return
				}

				slog.Error("Panic recovered", "details", string(logJSON))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}

// 判断是否需要记录请求体
func shouldLogRequestBody(method, contentType string) bool {
	if method != "POST" && method != "PUT" && method != "PATCH" {
		return false
	}

	// 只记录JSON和表单数据
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "application/x-www-form-urlencoded")
}

// 读取请求体
func readRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	// 重新设置请求体，供后续处理使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 限制日志长度，避免日志过大
	if len(body) > 1024 {
		return string(body[:1024]) + "...[truncated]"
	}

	return string(body)
}
