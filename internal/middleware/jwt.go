package middleware

import (
	"net/http"
	"path"
	"strings"
	"time"

	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

// JWT 白名单路径（精确匹配和前缀匹配）
var (
	// 精确匹配的路径
	jwtExactWhitelist = map[string]bool{
		"/api/login":        true,
		"/api/refreshToken": true,
		"/health":           true,
		"/health/ready":     true,
		"/health/live":      true,
		"/metrics":          true,
	}
	// 前缀匹配的路径
	jwtPrefixWhitelist = []string{
		"/api/v1/upload/images",
		"/console",
		"/static",
	}
)

// isPathWhitelisted 检查路径是否在白名单中（规范化后匹配）
func isPathWhitelisted(reqPath string) bool {
	// 规范化路径，防止 /api/login/../admin 类型的绕过
	cleanPath := path.Clean(reqPath)

	// 移除查询参数
	if idx := strings.Index(cleanPath, "?"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}

	// 精确匹配
	if jwtExactWhitelist[cleanPath] {
		return true
	}

	// 前缀匹配
	for _, prefix := range jwtPrefixWhitelist {
		if strings.HasPrefix(cleanPath, prefix) {
			return true
		}
	}

	return false
}

// JwtVerify JWT验证中间件（增强版）
func JwtVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 规范化路径检查白名单
		if isPathWhitelisted(c.Request.URL.Path) {
			c.Next()
			return
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

		// 检查 Token 是否在黑名单中（已注销）
		if jwttoken.IsTokenBlacklisted(tokenStr) {
			result.With(c, http.StatusUnauthorized, "Token已失效，请重新登录", nil)
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
