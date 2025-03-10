package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

func GinLogger() gin.HandlerFunc {
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	// type LogFormatter func(params LogFormatterParams) string 这里的LogFormatterParams是一个格式化日志参数的结构体
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		// 127.0.0.1 - [Sun, 22 Nov 2020 17:09:53 CST] "GET /ping HTTP/1.1 200 56.113µs "curl/7.64.1" "
		return fmt.Sprintf("%s - [%s] %s %d %s %13v %s\n",
			param.TimeStamp.Format("2006-01-02 15:01:05"), //请求时间
			param.ClientIP,     //请求客户端的IP地址
			param.Method,       //请求方法
			param.StatusCode,   //http响应码
			param.Latency,      //请求到响应的延时
			param.Path,         //路由路径
			param.ErrorMessage, //如果有错误,也打印错误信息
		)
	})
}

// GinRecovery recover掉项目可能出现的panic，并使用slog记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					slog.Errorf("path: %s,error: %s,request: %s ", c.Request.URL.Path, err, string(httpRequest))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					slog.Errorf("[Recovery from panic] err: %s;request:%s stack: %s", err, string(httpRequest), string(debug.Stack()))
				} else {
					slog.Errorf("[Recovery from panic] err: %s request: %s", err, string(httpRequest))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
