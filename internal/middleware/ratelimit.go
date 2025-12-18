package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitee.com/jieepre/go-site/internal/cache"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	MaxRequests int                       // 最大请求数
	Window      time.Duration             // 时间窗口
	KeyFunc     func(*gin.Context) string // 生成限流键的函数
}

// RateLimit 限流中间件
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成限流键
		key := config.KeyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

		// 获取当前请求计数
		currentCount, err := cache.Incr(rateLimitKey)
		if err != nil {
			// Redis 不可用时，使用内存限流器降级
			memLimiter := GetMemoryLimiter()
			if !memLimiter.Allow(key) {
				result.With(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试", nil)
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// 如果是第一次请求，设置过期时间
		if currentCount == 1 {
			_ = cache.Expire(rateLimitKey, config.Window)
		}

		// 检查是否超过限制
		if currentCount > int64(config.MaxRequests) {
			result.With(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试", nil)
			c.Abort()
			return
		}

		// 设置响应头，告知客户端限流状态
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(config.MaxRequests)-currentCount, 10))

		c.Next()
	}
}

// IPBasedRateLimit IP基础限流
func IPBasedRateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxRequests: maxRequests,
		Window:      window,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})
}

// UserBasedRateLimit 用户基础限流
func UserBasedRateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxRequests: maxRequests,
		Window:      window,
		KeyFunc: func(c *gin.Context) string {
			// 从JWT中获取用户名
			if username, exists := c.Get("username"); exists {
				return fmt.Sprintf("user:%s", username)
			}
			// 如果没有用户信息，使用IP
			return c.ClientIP()
		},
	})
}

// APIRateLimit API特定限流
func APIRateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxRequests: maxRequests,
		Window:      window,
		KeyFunc: func(c *gin.Context) string {
			// 组合IP和路径作为键
			path := c.Request.URL.Path
			ip := c.ClientIP()
			return fmt.Sprintf("api:%s:%s", ip, path)
		},
	})
}

// LoginRateLimit 登录限流（更严格）
func LoginRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxRequests: 5, // 5分钟内最多5次登录尝试
		Window:      5 * time.Minute,
		KeyFunc: func(c *gin.Context) string {
			ip := c.ClientIP()
			path := c.Request.URL.Path
			if strings.Contains(path, "login") {
				return fmt.Sprintf("login:%s", ip)
			}
			return ""
		},
	})
}
