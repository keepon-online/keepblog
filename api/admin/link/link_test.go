package link

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/internal/middleware"
	"github.com/gin-gonic/gin"
)

// 3.6.1 SSRF 修复的回归测试：友链检测接口的拨号层校验
// （回环/私网/链路本地/未指定/组播一律拒绝），防止被用作内网探测器。

// Control 校验只做 IP 分类判断，不产生网络 IO，可离线穷举各地址类别
func TestSafeDialerControlBlocksNonPublicIPs(t *testing.T) {
	blocked := []string{
		"127.0.0.1:80",       // IPv4 回环
		"10.0.12.5:8080",     // RFC1918 私网
		"172.16.0.1:80",      // RFC1918 私网
		"192.168.1.1:443",    // RFC1918 私网
		"169.254.169.254:80", // 链路本地（云元数据服务）
		"0.0.0.0:80",         // 未指定地址
		"224.0.0.1:80",       // 组播
		"[::1]:80",           // IPv6 回环
		"[fc00::1]:80",       // IPv6 私网
		"[fe80::1]:80",       // IPv6 链路本地
		"not-an-ip:80",       // 无法解析
	}
	for _, addr := range blocked {
		if err := safeDialer.Control("tcp", addr, nil); err == nil {
			t.Errorf("Control(%q) 应拒绝", addr)
		}
	}

	if err := safeDialer.Control("tcp", "8.8.8.8:443", nil); err != nil {
		t.Errorf("公网地址不应被拒: %v", err)
	}
	if err := safeDialer.Control("tcp", "[2606:4700::1111]:443", nil); err != nil {
		t.Errorf("公网 IPv6 地址不应被拒: %v", err)
	}
}

// CheckLink 处理器：内网目标与非 http(s) scheme 必须以“不可达/无效”收场，
// 不区分是拨号被拒还是真的超时——对调用方两者都是不可访问。
func TestCheckLinkRejectsIntranetTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// 空 url 分支走 result.Error，需要 ErrorHandler 统一渲染
	r.Use(middleware.ErrorHandler())
	h := &Handler{}
	r.GET("/api/v1/site/link/check", h.CheckLink)

	cases := []struct {
		name string
		url  string
		want string
	}{
		{"loopback", "http://127.0.0.1:8080/x", "访问超时或无法连接"},
		{"private", "http://192.168.1.1/admin", "访问超时或无法连接"},
		{"metadata", "http://169.254.169.254/latest/meta-data/", "访问超时或无法连接"},
		{"bad_scheme", "file:///etc/passwd", "无效链接"},
		{"ftp_scheme", "ftp://example.com/pub", "无效链接"},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/site/link/check?url="+c.url, nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("%s: 状态码 = %d, 期望 200（检测结果在 payload 中）", c.name, w.Code)
		}
		if !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("%s: 应包含 %q: %s", c.name, c.want, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), `"status":0`) {
			t.Errorf("%s: 检测结果应为 status=0: %s", c.name, w.Body.String())
		}
	}

	// 空 url 属于参数错误，走 400 而非检测结果
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/site/link/check?url=", nil))
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "url 不能为空") {
		t.Errorf("空 url 应 400 + 明确错误: %d %s", w.Code, w.Body.String())
	}
}
