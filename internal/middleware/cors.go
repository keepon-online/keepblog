package middleware

import (
	"os"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
)

// 允许的域名白名单（从环境变量读取，逗号分隔）
var allowedOrigins = initAllowedOrigins()

func initAllowedOrigins() map[string]bool {
	origins := map[string]bool{
		"http://localhost:8589":     true,
		"http://localhost:3000":     true,
		"http://127.0.0.1:8589":     true,
		"https://www.keepon.online": true,
		"https://keepon.online":     true,
	}

	// 从环境变量追加
	envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if envOrigins != "" {
		for _, origin := range strings.Split(envOrigins, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				origins[origin] = true
			}
		}
	}

	return origins
}

// isOriginAllowed 检查 Origin 是否在白名单中
func isOriginAllowed(origin string) bool {
	return allowedOrigins[origin]
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			// 检查 Origin 是否在白名单中
			if isOriginAllowed(origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, X-CSRF-Token")
				c.Header("Access-Control-Expose-Headers", "Content-Length, X-Token-Refresh")
				c.Header("Access-Control-Max-Age", "86400") // 24小时
				c.Header("Access-Control-Allow-Credentials", "true")
			} else {
				slog.Warnf("CORS blocked origin: %s", origin)
			}
		}

		if method == "OPTIONS" {
			if origin != "" && isOriginAllowed(origin) {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}

		c.Next()
	}
}
