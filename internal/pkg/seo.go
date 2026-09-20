package pkg

import (
	"bytes"
	"encoding/json"
	"gitee.com/jieepre/keepblog/internal/model/system"
	"github.com/gookit/slog"
	"io"
	"net/http"
	"strings"
)

// CanonicalURL 拼接站内绝对 URL：base 去掉尾部斜杠后接 path（path 以 / 开头）。
// 供各页面 handler 生成 canonical / og:url 使用。
func CanonicalURL(base, path string) string {
	return strings.TrimRight(base, "/") + path
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
