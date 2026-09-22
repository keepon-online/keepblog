package ai

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// mockUpstream 模拟 OpenAI 兼容的流式服务：校验请求头/体，按块吐出 delta。
func mockUpstream(t *testing.T, deltas []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %s, want /chat/completions", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("Authorization = %q", auth)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		for _, d := range deltas {
			fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"content\":%s}}]}\n\n", quoteJSON(d))
			flusher.Flush()
		}
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":7,\"total_tokens\":17}}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
}

func quoteJSON(s string) string {
	return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""
}

func TestStreamChat_DeltasAndUsage(t *testing.T) {
	srv := mockUpstream(t, []string{"你", "好", "，世界"})
	defer srv.Close()

	c := New(srv.URL, "test-key", "test-model", 10)
	var sb strings.Builder
	usage, err := c.StreamChat(context.Background(), StreamRequest{
		Messages:  []Message{{Role: "user", Content: "hi"}},
		MaxTokens: 100,
		OnDelta:   func(delta string) { sb.WriteString(delta) },
	})
	if err != nil {
		t.Fatalf("StreamChat: %v", err)
	}
	if got := sb.String(); got != "你好，世界" {
		t.Errorf("拼接结果 = %q, want %q", got, "你好，世界")
	}
	if usage.TotalTokens != 17 || usage.PromptTokens != 10 {
		t.Errorf("usage = %+v", usage)
	}
}

func TestStreamChat_UpstreamError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"error":"bad key"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "k", "m", 10)
	_, err := c.StreamChat(context.Background(), StreamRequest{Messages: []Message{{Role: "user", Content: "x"}}})
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Errorf("err = %v, want 包含 401", err)
	}
}

func TestStreamChat_ContextCancel(t *testing.T) {
	block := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block // 挂住不响应，等客户端取消
	}))
	defer srv.Close()
	defer close(block)

	c := New(srv.URL, "k", "m", 30)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err := c.StreamChat(ctx, StreamRequest{Messages: []Message{{Role: "user", Content: "x"}}})
	if err == nil {
		t.Error("取消后应返回错误")
	}
}

func TestStreamChat_SkipsMalformedChunk(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		f := w.(http.Flusher)
		fmt.Fprint(w, "data: not-json\n\n")
		fmt.Fprint(w, ": keep-alive comment\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
		f.Flush()
	}))
	defer srv.Close()

	c := New(srv.URL, "k", "m", 10)
	var got string
	_, err := c.StreamChat(context.Background(), StreamRequest{
		Messages: []Message{{Role: "user", Content: "x"}},
		OnDelta:  func(d string) { got += d },
	})
	if err != nil {
		t.Fatalf("StreamChat: %v", err)
	}
	if got != "ok" {
		t.Errorf("got = %q, want ok（坏块应跳过）", got)
	}
}

// TestStreamChat_ReasoningCallback 推理模型思考内容单独回调、不混入正文。
func TestStreamChat_ReasoningCallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		f := w.(http.Flusher)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"reasoning_content\":\"想一想\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"reasoning_content\":\"再想想\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"答案\"}}]}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
		f.Flush()
	}))
	defer srv.Close()

	var content, reasoning string
	_, err := New(srv.URL, "k", "m", 10).StreamChat(context.Background(), StreamRequest{
		Messages:    []Message{{Role: "user", Content: "x"}},
		OnDelta:     func(d string) { content += d },
		OnReasoning: func(d string) { reasoning += d },
	})
	if err != nil {
		t.Fatal(err)
	}
	if reasoning != "想一想再想想" {
		t.Errorf("reasoning = %q", reasoning)
	}
	if content != "答案" {
		t.Errorf("content = %q（思考不得混入正文）", content)
	}
}
