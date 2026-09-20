package templates

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/model/system"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
)

// parseTemplates 以应用启动同款方式解析全部模板。
// 模板语法错误会让 routers.go 的 template.Must 在启动时 panic，
// 这里把它提前到测试期暴露。
func parseTemplates(t *testing.T) *template.Template {
	t.Helper()
	return template.Must(template.New("").Funcs(core.TemplateFunc()).ParseFS(Fs, "**/*.html"))
}

func TestTemplatesParse(t *testing.T) {
	if len(parseTemplates(t).Templates()) == 0 {
		t.Fatal("no templates parsed")
	}
}

func siteFixture() *system.WebSite {
	return &system.WebSite{
		Title:       "测试站点",
		Description: "站点描述",
		Keywords:    "关键词",
		URL:         "https://example.com",
		BaiduSite:   "baidu-code",
		GoogleSite:  "google-code",
		BingSite:    "bing-code",
	}
}

// TestHeadPostMeta 渲染文章页 head：canonical、差异化 title 与合法的
// BlogPosting JSON-LD 都在此产出，任何一处模板写错都会破坏 SEO 元信息。
func TestHeadPostMeta(t *testing.T) {
	tpl := parseTemplates(t)
	data := map[string]any{
		"site": siteFixture(),
		"posts": &model.Post{
			Title:            "文章标题",
			PostSlug:         "abc123",
			Summary:          "文章摘要 \"带引号\"",
			Author:           "作者",
			CoverImage:       "https://example.com/cover.png",
			WordCount:        1234,
			CreateTime:       1700000000,
			LastModifiedTime: 1700003600,
		},
		"canonical": "https://example.com/post/abc123",
	}

	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "layout/head.html", data); err != nil {
		t.Fatalf("render head: %v", err)
	}
	out := b.String()

	if !strings.Contains(out, "<title>文章标题 | 测试站点</title>") {
		t.Errorf("title 不符合预期, got:\n%s", out)
	}
	if !strings.Contains(out, `rel="canonical"`) || !strings.Contains(out, "https://example.com/post/abc123") {
		t.Errorf("missing canonical, got:\n%s", out)
	}
	if !strings.Contains(out, "property=\"article:published_time\"") {
		t.Errorf("missing article:published_time")
	}

	// 抽出 JSON-LD 并验证是合法 JSON（模板插值经 JS 上下文转义后仍须可解析）。
	start := strings.Index(out, `<script type="application/ld+json">`)
	end := strings.Index(out, "</script>")
	if start < 0 || end < start {
		t.Fatalf("missing JSON-LD block, got:\n%s", out)
	}
	raw := out[start+len(`<script type="application/ld+json">`) : end]
	var ld map[string]any
	if err := json.Unmarshal([]byte(raw), &ld); err != nil {
		t.Fatalf("JSON-LD 不是合法 JSON: %v\n%s", err, raw)
	}
	if ld["@type"] != "BlogPosting" {
		t.Errorf("JSON-LD @type = %v, want BlogPosting", ld["@type"])
	}
	if ld["headline"] != "文章标题" {
		t.Errorf("JSON-LD headline = %v", ld["headline"])
	}
}

// TestHeadListMeta 渲染列表页 head：title 带页面名、description 用 pagedesc
// 覆盖全站默认值。
func TestHeadListMeta(t *testing.T) {
	tpl := parseTemplates(t)
	data := map[string]any{
		"site":      siteFixture(),
		"title":     "标签：Go",
		"pagedesc":  "标签「Go」下的全部文章，共 3 篇。",
		"canonical": "https://example.com/tags/Go",
	}

	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "layout/head.html", data); err != nil {
		t.Fatalf("render head: %v", err)
	}
	out := b.String()

	if !strings.Contains(out, "<title>标签：Go | 测试站点</title>") {
		t.Errorf("title 不符合预期, got:\n%s", out)
	}
	if !strings.Contains(out, `<meta content="标签「Go」下的全部文章，共 3 篇。"`) {
		t.Errorf("description 未使用 pagedesc, got:\n%s", out)
	}
	if !strings.Contains(out, "https://example.com/tags/Go") {
		t.Errorf("missing canonical url")
	}
}

// TestHeadVerificationMetas 三家站长平台的验证 meta 都应在配置非空时输出。
func TestHeadVerificationMetas(t *testing.T) {
	tpl := parseTemplates(t)
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "layout/head.html", map[string]any{"site": siteFixture()}); err != nil {
		t.Fatalf("render head: %v", err)
	}
	out := b.String()

	for _, want := range []string{
		`name="baidu-site-verification"`,
		`content="google-code"`,
		`name="google-site-verification"`,
		`content="bing-code"`,
		`name="msvalidate.01"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s, got:\n%s", want, out)
		}
	}
}
