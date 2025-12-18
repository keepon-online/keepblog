package middleware

import (
	"fmt"
	"net/http"
	"time"

	"gitee.com/jieepre/go-site/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Success   bool                   `json:"success"`
	Code      int                    `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Path      string                 `json:"path"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) == 0 {
			return
		}

		// 获取最后一个错误
		err := c.Errors.Last().Err

		// 获取请求ID
		requestID := ""
		if id, exists := c.Get("request_id"); exists {
			requestID = fmt.Sprintf("%v", id)
		}

		// 处理不同类型的错误
		var appErr *errors.AppError
		var httpStatus int
		var errorResponse ErrorResponse

		switch e := err.(type) {
		case *errors.AppError:
			// 应用自定义错误
			appErr = e
			if requestID != "" {
				appErr = appErr.WithRequestID(requestID)
			}

			httpStatus = appErr.HTTPStatus()
			errorResponse = ErrorResponse{
				Success:   false,
				Code:      int(appErr.Code),
				Message:   appErr.Message,
				Details:   appErr.Details,
				RequestID: appErr.RequestID,
				Timestamp: appErr.Timestamp.Format(time.DateTime),
				Path:      c.Request.URL.Path,
				Metadata:  appErr.Metadata,
			}

			// 记录错误日志
			logError(appErr, c)

		default:
			// 其他类型错误，包装为内部错误
			appErr = errors.Internal("服务器内部错误").WithRequestID(requestID)
			httpStatus = http.StatusInternalServerError
			errorResponse = ErrorResponse{
				Success:   false,
				Code:      int(errors.ErrInternal),
				Message:   "服务器内部错误",
				RequestID: requestID,
				Timestamp: appErr.Timestamp.Format(time.DateTime),
				Path:      c.Request.URL.Path,
			}

			// 记录原始错误
			slog.WithFields(map[string]interface{}{
				"request_id":     requestID,
				"path":           c.Request.URL.Path,
				"method":         c.Request.Method,
				"client_ip":      c.ClientIP(),
				"original_error": err.Error(),
			}).Error("Unhandled error occurred")
		}

		// 如果响应已经写入，则无法修改
		if c.Writer.Written() {
			return
		}

		// 设置响应
		c.JSON(httpStatus, errorResponse)
		c.Abort()
	}
}

// logError 记录错误日志
func logError(err *errors.AppError, c *gin.Context) {
	fields := map[string]interface{}{
		"error_code": err.Code,
		"error_msg":  err.Message,
		"request_id": err.RequestID,
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
		"client_ip":  c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
	}

	if err.Details != "" {
		fields["details"] = err.Details
	}

	if err.Metadata != nil {
		fields["metadata"] = err.Metadata
	}

	// 根据错误级别选择日志级别
	switch {
	case err.Code >= 1000 && err.Code < 2000: // 系统错误
		if err.Stack != "" {
			fields["stack"] = err.Stack
		}
		slog.WithFields(fields).Error("System error occurred")

	case err.Code >= 2000 && err.Code < 3000: // 认证错误
		slog.WithFields(fields).Warn("Authentication error occurred")

	case err.Code >= 3000 && err.Code < 4000: // 参数错误
		slog.WithFields(fields).Info("Parameter validation error occurred")

	case err.Code >= 4000 && err.Code < 5000: // 业务错误
		slog.WithFields(fields).Info("Business logic error occurred")

	case err.Code >= 5000: // 外部服务错误
		slog.WithFields(fields).Error("External service error occurred")

	default:
		slog.WithFields(fields).Error("Unknown error occurred")
	}
}

// AbortWithError 中断请求并返回错误
func AbortWithError(c *gin.Context, err *errors.AppError) {
	c.Error(err)
	c.Abort()
}

// AbortWithAppError 便捷方法：中断请求并返回应用错误
func AbortWithAppError(c *gin.Context, code errors.ErrorCode, message string) {
	err := errors.New(code, message)
	if requestID, exists := c.Get("request_id"); exists {
		err = err.WithRequestID(fmt.Sprintf("%v", requestID))
	}
	AbortWithError(c, err)
}

// HandleDBError 处理数据库错误
func HandleDBError(c *gin.Context, err error) {
	dbErr := errors.Wrap(err, errors.ErrDatabase, "数据库操作失败")
	if requestID, exists := c.Get("request_id"); exists {
		dbErr = dbErr.WithRequestID(fmt.Sprintf("%v", requestID))
	}
	AbortWithError(c, dbErr)
}

// HandleValidationError 处理参数验证错误
func HandleValidationError(c *gin.Context, field string, message string) {
	err := errors.InvalidParam(fmt.Sprintf("字段 '%s' %s", field, message))
	if requestID, exists := c.Get("request_id"); exists {
		err = err.WithRequestID(fmt.Sprintf("%v", requestID))
	}
	AbortWithError(c, err)
}
