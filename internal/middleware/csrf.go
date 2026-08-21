package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
)

// CSRFConfig CSRF 配置
type CSRFConfig struct {
	TokenLength   int           // Token 长度
	CookieName    string        // Cookie 名称
	HeaderName    string        // 请求头名称
	CookieMaxAge  time.Duration // Cookie 有效期
	Secure        bool          // 仅 HTTPS
	SameSite      http.SameSite // SameSite 策略
	IgnoreMethods []string      // 忽略的 HTTP 方法
}

// DefaultCSRFConfig 默认 CSRF 配置
var DefaultCSRFConfig = CSRFConfig{
	TokenLength:   32,
	CookieName:    "csrf_token",
	HeaderName:    "X-CSRF-Token",
	CookieMaxAge:  24 * time.Hour,
	Secure:        true, // 生产环境必须为 true，强制 HTTPS
	SameSite:      http.SameSiteStrictMode,
	IgnoreMethods: []string{"GET", "HEAD", "OPTIONS"},
}

// tokenStore 存储已使用的 token（防止重放攻击）
var (
	usedTokens     = make(map[string]time.Time)
	usedTokensMu   sync.RWMutex
	tokenCleanupOn sync.Once
)

// generateToken 生成随机 CSRF token
func generateCSRFToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// isMethodIgnored 检查方法是否应该忽略 CSRF 验证
func isMethodIgnored(method string, ignoreMethods []string) bool {
	for _, m := range ignoreMethods {
		if m == method {
			return true
		}
	}
	return false
}

// cleanupExpiredTokens 定期清理过期 token
func cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		usedTokensMu.Lock()
		now := time.Now()
		for token, expiry := range usedTokens {
			if now.After(expiry) {
				delete(usedTokens, token)
			}
		}
		usedTokensMu.Unlock()
	}
}

// CSRF 返回 CSRF 保护中间件
func CSRF() gin.HandlerFunc {
	return CSRFWithConfig(DefaultCSRFConfig)
}

// CSRFWithConfig 返回带配置的 CSRF 保护中间件
func CSRFWithConfig(config CSRFConfig) gin.HandlerFunc {
	// 启动清理协程
	tokenCleanupOn.Do(func() {
		go cleanupExpiredTokens()
	})

	return func(c *gin.Context) {
		// 获取或生成 CSRF token
		token, err := c.Cookie(config.CookieName)
		if err != nil || token == "" {
			token, err = generateCSRFToken(config.TokenLength)
			if err != nil {
				result.Error(c, "Failed to generate CSRF token")
				c.Abort()
				return
			}

			c.SetCookie(
				config.CookieName,
				token,
				int(config.CookieMaxAge.Seconds()),
				"/",
				"",
				config.Secure,
				true, // HttpOnly
			)
		}

		// 将 token 存入上下文供模板使用
		c.Set("csrf_token", token)

		// 忽略安全方法
		if isMethodIgnored(c.Request.Method, config.IgnoreMethods) {
			c.Next()
			return
		}

		// 验证请求中的 token
		requestToken := c.GetHeader(config.HeaderName)
		if requestToken == "" {
			requestToken = c.PostForm("csrf_token")
		}

		if requestToken == "" {
			result.With(c, http.StatusForbidden, "CSRF token missing", nil)
			c.Abort()
			return
		}

		if requestToken != token {
			result.With(c, http.StatusForbidden, "CSRF token mismatch", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
