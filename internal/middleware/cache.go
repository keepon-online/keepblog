package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitee.com/jieepre/keepblog/internal/cache"

	"github.com/gin-gonic/gin"
)

// CacheMiddleware 缓存中间件
func CacheMiddleware(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只缓存GET请求
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// 生成缓存键
		cacheKey := generateCacheKey(c)

		// 尝试从缓存获取
		var cachedResponse CachedResponse
		if cache.Exists(cacheKey) {
			if err := cache.Get(cacheKey, &cachedResponse); err == nil {
				// 设置响应头
				for key, value := range cachedResponse.Headers {
					c.Header(key, value)
				}
				c.Data(cachedResponse.StatusCode, cachedResponse.ContentType, cachedResponse.Body)
				c.Abort()
				return
			}
		}

		// 创建响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		c.Next()

		// 缓存响应
		if writer.status >= 200 && writer.status < 300 {
			headers := make(map[string]string)
			for key, values := range writer.Header() {
				if len(values) > 0 {
					headers[key] = values[0]
				}
			}

			cachedResp := CachedResponse{
				StatusCode:  writer.status,
				Headers:     headers,
				Body:        writer.body.Bytes(),
				ContentType: writer.Header().Get("Content-Type"),
			}

			_ = cache.Set(cacheKey, cachedResp, duration)
		}
	}
}

// CachedResponse 缓存的响应
type CachedResponse struct {
	StatusCode  int               `json:"status_code"`
	Headers     map[string]string `json:"headers"`
	Body        []byte            `json:"body"`
	ContentType string            `json:"content_type"`
}

// responseWriter 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// generateCacheKey 生成缓存键
func generateCacheKey(c *gin.Context) string {
	// 使用请求路径和查询参数生成缓存键
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	// 排除一些不需要缓存的路径
	excludePaths := []string{"/api/", "/admin/"}
	for _, excludePath := range excludePaths {
		if strings.Contains(path, excludePath) {
			return ""
		}
	}

	key := fmt.Sprintf("cache:%s", path)
	if query != "" {
		key = fmt.Sprintf("%s?%s", key, query)
	}

	// 使用MD5哈希来缩短键名
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("page:%x", hash)
}

// ClearCacheByPattern 清除匹配模式的缓存
func ClearCacheByPattern(pattern string) error {
	return cache.DeletePattern(pattern)
}
