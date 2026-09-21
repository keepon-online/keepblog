// Package notify 文章发布通知：发布动作触发后异步推送 Webhook / Telegram，
// 失败只记日志，绝不阻塞或影响发布本身。
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/model"
	"github.com/gookit/slog"
)

// httpClient 复用连接；通知是低频操作，超时收紧避免 goroutine 堆积。
var httpClient = &http.Client{Timeout: 5 * time.Second}

// Payload 推送到 Webhook 的 JSON 结构。
type Payload struct {
	Event     string `json:"event"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Author    string `json:"author,omitempty"`
	Summary   string `json:"summary,omitempty"`
	PubTime   int64  `json:"pubTime"`
	Scheduled bool   `json:"scheduled"`
}

// PostPublished 文章发布通知入口。异步执行，调用方无需等待。
// baseURL 为站点地址（带尾斜亦可），slug 由内部拼接。
func PostPublished(post model.Post, baseURL string) {
	cfg := config.Get().Notify
	if cfg == nil || (cfg.Webhook == "" && cfg.Telegram.Token == "") {
		return
	}
	url := strings.TrimRight(baseURL, "/") + "/post/" + post.PostSlug
	p := Payload{
		Event:   "post.published",
		Title:   post.Title,
		URL:     url,
		Author:  post.Author,
		Summary: post.Summary,
		PubTime: int64(post.PubTime),
		// 定时发布：pub_time 在未来，通知先行到达，标注到点时间
		Scheduled: post.PubTime > uint64(time.Now().Unix()),
	}
	go func() {
		if cfg.Webhook != "" {
			sendWebhook(cfg.Webhook, p)
		}
		if cfg.Telegram.Token != "" && cfg.Telegram.ChatId != "" {
			sendTelegram(cfg.Telegram.Token, cfg.Telegram.ChatId, p)
		}
	}()
}

func sendWebhook(webhookURL string, p Payload) {
	body, err := json.Marshal(p)
	if err != nil {
		slog.Errorf("[notify] webhook 序列化失败: %s", err.Error())
		return
	}
	resp, err := httpClient.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		slog.Errorf("[notify] webhook 推送失败: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		slog.Errorf("[notify] webhook 返回异常状态码 %d", resp.StatusCode)
	}
}

func sendTelegram(token, chatId string, p Payload) {
	text := fmt.Sprintf("📰 新文章：%s\n%s", p.Title, p.URL)
	if p.Scheduled {
		text = fmt.Sprintf("⏰ 定时发布：%s（%s 上线）\n%s", p.Title,
			time.Unix(p.PubTime, 0).Format("2006-01-02 15:04"), p.URL)
	}
	body, _ := json.Marshal(map[string]string{
		"chat_id": chatId,
		"text":    text,
	})
	resp, err := httpClient.Post("https://api.telegram.org/bot"+token+"/sendMessage",
		"application/json", bytes.NewReader(body))
	if err != nil {
		slog.Errorf("[notify] telegram 推送失败: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		slog.Errorf("[notify] telegram 返回异常状态码 %d", resp.StatusCode)
	}
}
