package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"net/http"

	"gitee.com/jieepre/keepblog/static"
	"github.com/gin-gonic/gin"
)

// 静态资源缓存策略：内容经 embed.FS 编译进二进制、随发版才变化，但没有
// 版本化 URL，不能使用 immutable 长强缓存——那样客户端会把旧文件钉死到
// 过期为止，发版后的 JS/CSS 修复永远到不了老访客（autoplay 修复上线后
// 浏览器仍报旧错的根因）。改为 no-cache + 内容哈希 ETag：浏览器每次回源
// 校验，内容未变回 304（无 body，开销极小），内容一变立即生效。
const staticCacheControl = "public, no-cache"

// StaticCacheMiddleware 为静态资源响应注入协商缓存头。
//
// 注意 HTTP 响应 header 必须在 body 写出前设置：静态文件由内层
// http.FileServer 写出，因此本中间件在 c.Next() 之前预设 header，
// 由 FileServer 在写出 2xx/304 时沿用（If-None-Match 命中时由
// ServeContent 直接写 304）。这与 gzip 中间件兼容（gzip 只压缩
// body，不改 header 语义）。
//
// 对极少数 404（资源不存在）也会带上缓存头，影响可忽略：浏览器对
// 不存在的资源本就不会重复强缓存命中，且静态资源路径稳定。
func StaticCacheMiddleware() gin.HandlerFunc {
	etag := staticEtag()
	return func(c *gin.Context) {
		// 只对静态资源的 GET/HEAD 生效
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			c.Header("Cache-Control", staticCacheControl)
			c.Header("ETag", etag)
		}
		c.Next()
	}
}

// staticEtag 对全部内嵌静态资源的路径与内容求哈希，作为统一的 ETag：
// 二进制不变则值稳定（304 可命中），任一资源内容变化即产生新值。
// 静态资源总量仅数 MB，启动时一次计算耗时可忽略。
func staticEtag() string {
	h := sha256.New()
	_ = fs.WalkDir(static.Static, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		h.Write([]byte(path))
		data, err := static.Static.ReadFile(path)
		if err != nil {
			return err
		}
		h.Write(data)
		return nil
	})
	// favicon 在独立的 embed FS 中，单独计入
	h.Write([]byte("favicon.ico"))
	if data, err := static.Favicon.ReadFile("favicon.ico"); err == nil {
		h.Write(data)
	}
	sum := hex.EncodeToString(h.Sum(nil))
	return `"` + sum[:16] + `"`
}
