package middleware

import (
	"net/http"
	"strings"
	"time"

	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

// JwtVerify JWT验证中间件（增强版）
func JwtVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 白名单路径，不需要验证
		whitelist := []string{
			"/api/login",
			"/api/refreshToken",
			"/api/v1/upload/images",
			"/console",
			"/health", // 健康检查
		}

		reqURI := c.Request.RequestURI
		for _, whitelistPath := range whitelist {
			if reqURI == whitelistPath || strings.HasPrefix(reqURI, whitelistPath) {
				c.Next()
				return
			}
		}

		// 获取Authorization头
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			result.With(c, http.StatusUnauthorized, "未提供认证信息", nil)
			c.Abort()
			return
		}

		// 检查Bearer格式
		if !strings.HasPrefix(authorization, "Bearer ") {
			result.With(c, http.StatusUnauthorized, "Token格式错误", nil)
			c.Abort()
			return
		}

		// 提取token
		tokenStr := strings.TrimPrefix(authorization, "Bearer ")
		if tokenStr == "" {
			result.With(c, http.StatusUnauthorized, "Token不能为空", nil)
			c.Abort()
			return
		}

		// 验证token
		claims, err := jwttoken.ParseToken(tokenStr)
		if err != nil {
			var message string

			switch {
			case strings.Contains(err.Error(), "expired"):
				message = "Token已过期"
			case strings.Contains(err.Error(), "not active yet"):
				message = "Token尚未生效"
			default:
				message = "Token验证失败"
			}

			result.With(c, http.StatusUnauthorized, message, nil)
			c.Abort()
			return
		}

		// 检查token是否即将过期（剩余时间小于30分钟）
		if time.Until(claims.ExpiresAt.Time) < 30*time.Minute {
			c.Header("X-Token-Refresh", "true") // 提示前端刷新token
		}

		// 将用户信息存入上下文
		c.Set("username", claims.Username)

		c.Next()
	}
}
