package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// 3.6.1 限流 fail-open 修复的回归测试：Redis 未启用时（cache.Enable=false，
// 测试进程的默认状态）必须降级到进程内固定窗口限流，而不是恒放行。

// 固定窗口计数器：窗口内超限拒绝、窗口到期自动重置
func TestWindowLimiterAllow(t *testing.T) {
	const key = "test-window-limiter-key"
	window := 60 * time.Millisecond

	for i := 0; i < 3; i++ {
		if !globalWindowLimiter.Allow(key, 3, window) {
			t.Fatalf("第 %d 次请求（限 3 次）不应被拒", i+1)
		}
	}
	if globalWindowLimiter.Allow(key, 3, window) {
		t.Fatal("窗口内第 4 次请求必须被拒")
	}

	time.Sleep(80 * time.Millisecond)
	if !globalWindowLimiter.Allow(key, 3, window) {
		t.Fatal("窗口过期后计数应重置")
	}
}

// 中间件级降级：Redis 不可用时按本端点阈值拦请求（429 走业务码）
func TestRateLimitFallbackWithoutRedis(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	probeHits := 0
	// 独立 key + 短窗口，避免与其它测试及后续请求串扰
	r.Use(RateLimit(RateLimitConfig{
		MaxRequests: 2,
		Window:      100 * time.Millisecond,
		KeyFunc: func(c *gin.Context) string {
			return "test-fallback-key"
		},
	}))
	r.Any("/probe", func(c *gin.Context) {
		probeHits++
		c.JSON(http.StatusOK, nil)
	})

	hit := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/probe", nil))
		return w
	}

	for i := 0; i < 2; i++ {
		w := hit()
		if w.Code != http.StatusOK {
			t.Fatalf("前 2 次请求应放行（第 %d 次），状态码 = %d", i+1, w.Code)
		}
	}
	if probeHits != 2 {
		t.Fatalf("前 2 次请求都应到达业务处理器，实际 %d 次", probeHits)
	}

	w := hit()
	if probeHits != 2 {
		t.Fatal("第 3 次请求必须被限流中间件拦截，不得到达业务处理器")
	}
	if !strings.Contains(w.Body.String(), `"code":429`) {
		t.Fatalf("限流响应应携带业务码 429: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "请求过于频繁") {
		t.Fatalf("限流响应应携带提示: %s", w.Body.String())
	}
}

// 空 key（如非登录路径上的 LoginRateLimit）直接放行，不参与限流
func TestRateLimitEmptyKeyPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(LoginRateLimit())
	r.Any("/not-login", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/not-login", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("非登录路径第 %d 次请求不应被登录限流拦截", i+1)
		}
	}
}
