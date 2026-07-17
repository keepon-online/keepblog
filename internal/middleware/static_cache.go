package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 静态资源强缓存时长：一年。静态资源通过 embed.FS 编译进二进制，
// 内容随发版才变化，适合长期强缓存；immutable 告知浏览器内容不会变，
// 无需在过期前发起条件请求。
const staticCacheControl = "public, max-age=31536000, immutable"

// StaticCacheMiddleware 为静态资源响应注入长缓存头。
//
// 注意 HTTP 响应 header 必须在 body 写出前设置：静态文件由内层
// http.FileServer 写出，因此本中间件在 c.Next() 之前预设 header，
// 由 FileServer 在写出 2xx/304 时沿用。这与 gzip 中间件兼容
// （gzip 只压缩 body，不改 header 语义）。
//
// 对极少数 404（资源不存在）也会带上缓存头，影响可忽略：浏览器对
// 不存在的资源本就不会重复强缓存命中，且静态资源路径稳定。
func StaticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对静态资源的 GET/HEAD 生效
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			c.Header("Cache-Control", staticCacheControl)
		}
		c.Next()
	}
}
