package pkg

import (
	"bytes"
	"encoding/json"
	"gitee.com/jieepre/go-site/internal/model/system"
	"github.com/gookit/slog"
	"io"
	"net/http"
	"strings"
)

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
