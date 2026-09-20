package page

import (
	"strings"
	"testing"
)

// renderPage 断言辅助：渲染并返回 HTML。
func renderPage(t *testing.T, total, pageN int, temp string) string {
	t.Helper()
	html, err := HandleIndex(total, pageN, temp)
	if err != nil {
		t.Fatalf("HandleIndex(%d, %d, %q) 报错: %v", total, pageN, temp, err)
	}
	return html
}

func TestHandleIndexFirstPageLinksToBase(t *testing.T) {
	cases := []struct{ temp, wantBase string }{
		{"page", `href="/"`},
		{"tags/go/page", `href="/tags/go"`},
		{"categories/java/page", `href="/categories/java"`},
		{"archives/page", `href="/archives"`},
	}
	for _, tc := range cases {
		// 第 2 页上看"1"这个快捷链接
		html := renderPage(t, 25, 2, tc.temp)
		if !strings.Contains(html, tc.wantBase) {
			t.Errorf("temp=%q: 第 1 页应指向基路径 %s, got:\n%s", tc.temp, tc.wantBase, html)
		}
	}
}

func TestHandleIndexCurrentPageNotSelfLinked(t *testing.T) {
	html := renderPage(t, 25, 1, "page")
	// 当前页：无 href、带 aria-current，且不出现 /page/1
	if !strings.Contains(html, `<li class="page-number current"><a aria-current="page">1</a></li>`) {
		t.Errorf("当前页应为无 href 锚点 + aria-current, got:\n%s", html)
	}
	if strings.Contains(html, "/page/1") {
		t.Errorf("不应再生成 /page/1 链接, got:\n%s", html)
	}
}

func TestHandleIndexDisabledArrowsHaveNoHref(t *testing.T) {
	html := renderPage(t, 25, 1, "page")
	if strings.Contains(html, `href="#"`) {
		t.Errorf("禁用箭头不应输出 href=\"#\", got:\n%s", html)
	}
	if !strings.Contains(html, `class="disabled"`) {
		t.Errorf("首页上一页应为禁用态")
	}
	// 末页的下一页同样禁用
	html = renderPage(t, 25, 3, "page")
	if strings.Contains(html, `href="#"`) {
		t.Errorf("末页下一页不应输出 href=\"#\", got:\n%s", html)
	}
}

func TestHandleIndexNormalNavigation(t *testing.T) {
	html := renderPage(t, 25, 2, "page")
	for _, want := range []string{
		`<li><a href="/" data-pjax-content="true">1</a></li>`,
		`href="/page/3"`,
		`class="extend next"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("第 2 页缺少 %s, got:\n%s", want, html)
		}
	}
	// 上一页在第 2 页应指向首页基路径
	if !strings.Contains(html, `class="extend prev" href="/"`) {
		t.Errorf("第 2 页的上一页应指向 /, got:\n%s", html)
	}
}

func TestHandleIndexEllipsis(t *testing.T) {
	html := renderPage(t, 150, 8, "page") // 15 页，中间位置出现省略号
	if !strings.Contains(html, ">...</a></li>") || !strings.Contains(html, `href="/"`) {
		t.Errorf("多页省略号与首页快捷链接缺失, got:\n%s", html)
	}
	if strings.Contains(html, "/page/1<") {
		t.Errorf("首页快捷链接不应指向 /page/1")
	}
}

func TestHandleIndexOverflowPage(t *testing.T) {
	if _, err := HandleIndex(25, 99, "page"); err == nil {
		t.Errorf("超界页码应报错")
	}
	// 单页内容不输出分页
	if html := renderPage(t, 5, 1, "page"); html != "" {
		t.Errorf("单页不应输出分页 HTML, got %s", html)
	}
}
