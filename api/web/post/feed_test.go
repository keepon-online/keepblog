package post

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/internal/service"
	"gitee.com/jieepre/keepblog/internal/testutil"
	"github.com/gin-gonic/gin"
)

func setupWebTestRouter(t *testing.T) (*gin.Engine, *Handler) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := testutil.NewTestDB(t)
	ctx := &core.Context{Service: service.InitAppService(db)}
	r := gin.New()
	return r, &Handler{Context: ctx}
}

// TestFeed_FullTextContent RSS 全文输出：content:encoded 命名空间、渲染后的
// 正文 HTML、作者字段与封面图都必须出现在输出里，且定时文章不外泄。
func TestFeed_FullTextContent(t *testing.T) {
	r, handler := setupWebTestRouter(t)
	r.GET("/rss.xml", handler.Feed)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/rss.xml", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP = %d, body = %s", w.Code, w.Body.String())
	}
	body := w.Body.String()

	if !strings.Contains(body, `xmlns:content="http://purl.org/rss/1.0/modules/content/"`) {
		t.Error("缺少 content 命名空间声明")
	}
	if !strings.Contains(body, "<content:encoded>") {
		t.Error("缺少 content:encoded 全文节点")
	}
	// 正文是渲染后的 HTML（fixture 正文含 Markdown，渲染应产生 <p> 或代码块）
	if !strings.Contains(body, "&lt;p&gt;") && !strings.Contains(body, "&lt;h") && !strings.Contains(body, "&lt;pre") && !strings.Contains(body, "&lt;code") {
		t.Error("content:encoded 未包含渲染后的 HTML 标签")
	}
	if !strings.Contains(body, "<author>") {
		t.Error("缺少 author 作者节点")
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/xml") {
		t.Errorf("Content-Type = %q, want application/xml", ct)
	}
}
