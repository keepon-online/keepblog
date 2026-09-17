package middleware

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitee.com/jieepre/keepblog/static"
	"github.com/gin-gonic/gin"
)

// 按真实路由装配方式（middleware + StaticFS + embed FS）验证协商缓存行为。
func newStaticTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jsEmbed, _ := fs.Sub(static.Static, "js")
	group := r.Group("/", StaticCacheMiddleware())
	group.StaticFS("/js", http.FS(jsEmbed))
	return r
}

func TestStaticCacheEtagRevalidation(t *testing.T) {
	router := newStaticTestRouter()

	// 首次请求：200，no-cache + ETag
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, httptest.NewRequest(http.MethodGet, "/js/pjax-music.js", nil))
	if w1.Code != http.StatusOK {
		t.Fatalf("首次请求状态码 = %d, 期望 200", w1.Code)
	}
	etag := w1.Header().Get("ETag")
	if etag == "" {
		t.Fatal("缺少 ETag 响应头")
	}
	if cc := w1.Header().Get("Cache-Control"); cc != "public, no-cache" {
		t.Fatalf("Cache-Control = %q, 期望 public, no-cache", cc)
	}

	// 携带相同 ETag 再请求：304 且无 body
	req := httptest.NewRequest(http.MethodGet, "/js/pjax-music.js", nil)
	req.Header.Set("If-None-Match", etag)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	if w2.Code != http.StatusNotModified {
		t.Fatalf("If-None-Match 命中时状态码 = %d, 期望 304", w2.Code)
	}
	if w2.Body.Len() != 0 {
		t.Fatalf("304 响应不应有 body, 实际 %d 字节", w2.Body.Len())
	}

	// 携带不同 ETag（模拟新版本）：200 全量返回
	req2 := httptest.NewRequest(http.MethodGet, "/js/pjax-music.js", nil)
	req2.Header.Set("If-None-Match", `"stale-build"`)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req2)
	if w3.Code != http.StatusOK {
		t.Fatalf("ETag 不匹配时状态码 = %d, 期望 200", w3.Code)
	}

	// 同一进程内多次构造中间件，ETag 应一致（内容未变）
	w4 := httptest.NewRecorder()
	newStaticTestRouter().ServeHTTP(w4, httptest.NewRequest(http.MethodGet, "/js/pjax-music.js", nil))
	if w4.Header().Get("ETag") != etag {
		t.Fatal("内容未变时两次构造的 ETag 不一致")
	}
}
