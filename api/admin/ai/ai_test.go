package ai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/middleware"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/internal/service"
	"gitee.com/jieepre/keepblog/internal/testutil"
	"github.com/gin-gonic/gin"
)

// setupAI 构造挂好 AI 路由的测试引擎与 mock 上游。
// 注意：admin 分组的 JWT 中间件在路由器装配层（internal/app），本包直测
// handler 行为，鉴权由路由层覆盖。
func setupAI(t *testing.T, upstream *httptest.Server) (*gin.Engine, *Handler) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := testutil.NewTestDB(t)
	ctx := &core.Context{Service: service.InitAppService(db)}
	h := &Handler{Context: ctx}

	r := gin.New()
	r.Use(middleware.ErrorHandler()) // result.Error 依赖它渲染错误信封
	r.GET("/api/v1/ai/status", h.Status)
	r.POST("/api/v1/ai/edit", h.Edit)
	return r, h
}

func TestStatus_DisabledWithoutKey(t *testing.T) {
	r, _ := setupAI(t, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/ai/status", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"enabled":false`) {
		t.Fatalf("status = %d %s", w.Code, w.Body.String())
	}
}

func TestEdit_UnconfiguredRejected(t *testing.T) {
	r, _ := setupAI(t, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/edit", strings.NewReader(`{"task":"summary","digest":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	// 未配置/参数错走统一错误信封，不应挂起为 SSE
	if w.Code == 200 || !strings.Contains(w.Body.String(), "AI 未配置") {
		t.Fatalf("code = %d body = %s", w.Code, w.Body.String())
	}
}

func TestEdit_SSEStream(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		f := w.(http.Flusher)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"润\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"色后\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{},\"finish_reason\":\"stop\"}],\"usage\":{\"total_tokens\":9}}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
		f.Flush()
	}))
	defer upstream.Close()

	// 通过 config 快照注入配置：ai.apiKey/baseURL 指向 mock
	t.Setenv("KEEPBLOG_AI_APIKEY", "test-key")
	// viper 绑定的环境变量需要重新 Load 才进快照；测试进程内直接操作 config
	if err := config.StoreForTest(map[string]any{
		"ai.apiKey": "test-key", "ai.baseURL": upstream.URL, "ai.model": "m",
		"ai.dailyQuota": 5,
	}); err != nil {
		t.Fatalf("注入测试配置: %v", err)
	}

	r, _ := setupAI(t, upstream)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/edit", strings.NewReader(`{"task":"polish","selection":"原文"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q", ct)
	}
	for _, need := range []string{
		`"type":"delta"`, `"content":"润"`, `"content":"色后"`, `"type":"done"`, `"total_tokens":9`,
	} {
		if !strings.Contains(body, need) {
			t.Errorf("SSE 输出缺少 %s\nbody: %s", need, body)
		}
	}
}
