package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitee.com/jieepre/keepblog/internal/model/system"
	"github.com/gookit/slog"
	"io"
	"net/http"
	"strings"
	"time"
)

// CanonicalURL 拼接站内绝对 URL：base 去掉尾部斜杠后接 path（path 以 / 开头）。
// 供各页面 handler 生成 canonical / og:url 使用。
func CanonicalURL(base, path string) string {
	return strings.TrimRight(base, "/") + path
}

// PushIndexNow 向 IndexNow 兼容搜索引擎（Bing/Yandex 等）提交 URL 列表。
// host 是站点域名（协议部分需已剥离），key 为 IndexNow 密钥；
// endpoint 传空串时用官方共享入口。
// 协议规定响应 200（受理）或 202（待校验 key）为成功，其余状态码返回错误。
func PushIndexNow(urls []string, host, key, endpoint string) error {
	if len(urls) == 0 {
		return nil
	}
	if endpoint == "" {
		endpoint = "https://api.indexnow.org/indexnow"
	}

	body, err := json.Marshal(map[string]any{
		"host":    host,
		"key":     key,
		"urlList": urls,
	})
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("indexnow 推送失败: HTTP %d, %s", resp.StatusCode, string(respBody))
	}
	slog.Infof("indexnow 推送 %d 条 URL: HTTP %d", len(urls), resp.StatusCode)
	return nil
}

func PushSite(urls []string, api string) *system.PushSite {
	ch := &http.Client{}
	var site system.PushSite
	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBufferString(strings.Join(urls, "\n")))
	if err != nil {
		slog.Error(err)
		return nil
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := ch.Do(req)
	if err != nil {
		slog.Error(err)
		return nil
	}
	defer func() { _ = resp.Body.Close() }()
	slog.Infof("sate:%d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err)
		return nil
	}
	slog.Infof("推送结果:%v", string(body))
	err = json.Unmarshal(body, &site)
	if err != nil {
		slog.Error(err)
		return nil
	}
	return &site
}
