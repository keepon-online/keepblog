package md

import (
	"strings"
	"testing"
)

// 危险链接地址（javascript: 等）必须以各种编码形态被拦下：
// 实体/数字实体/大小写都会在浏览器侧还原成可执行 scheme。
func TestLinkDestinationSanitization(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"plain", `[x](javascript:alert(1))`},
		{"entity_colon", `[x](javascript&colon;alert(1))`},
		{"numeric_entity", `[x](&#x6a;avascript:alert(1))`},
		{"mixed_case", `[x](JaVaScRiPt:alert(1))`},
		{"tab_inside", "[x](java\tscript:alert(1))"},
		{"image_entity", `![i](javascript&colon;alert(1))`},
	}
	for _, c := range cases {
		out := ToHTML([]byte(c.in))
		if strings.Contains(strings.ToLower(out), "javascript:alert") {
			t.Errorf("%s: dangerous URL survived: %s", c.name, out)
		}
		t.Logf("%-14s => %s", c.name, out)
	}
}

// 正常链接必须保持可用，不能被误伤。
func TestLinkDestinationPreserved(t *testing.T) {
	ok := []struct {
		in   string
		want string
	}{
		{`[x](https://example.com/a?b=1&c=2)`, `href="https://example.com/a?b=1&amp;c=2"`},
		{`[x](#section)`, `href="#section"`},
		{`[x](mailto:a@b.com)`, `href="mailto:a@b.com"`},
		{`[x](/post/abc)`, `href="/post/abc"`},
	}
	for _, c := range ok {
		out := ToHTML([]byte(c.in))
		if !strings.Contains(out, c.want) {
			t.Errorf("%s: want %q in %q", c.in, c.want, out)
		}
	}
}