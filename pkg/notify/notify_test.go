package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestSendWebhook_Payload 验证 webhook 推送的 JSON 载荷与状态码处理。
func TestSendWebhook_Payload(t *testing.T) {
	var got Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("解码失败: %v", err)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	p := Payload{
		Event: "post.published", Title: "测试文章", URL: "https://x.dev/post/abc",
		Author: "tester", PubTime: time.Now().Add(time.Hour).Unix(), Scheduled: true,
	}
	sendWebhook(srv.URL, p)

	if got.Event != "post.published" || got.Title != "测试文章" || got.URL != "https://x.dev/post/abc" {
		t.Errorf("载荷字段错误: %+v", got)
	}
	if !got.Scheduled {
		t.Error("定时发布标识未传递")
	}
}

// TestSendWebhook_ServerError 服务端异常只记日志不 panic。
func TestSendWebhook_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()
	sendWebhook(srv.URL, Payload{Event: "post.published", Title: "t"})
}
