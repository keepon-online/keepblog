package md

import (
	"strings"
	"testing"
)

// 含标题、列表、代码块、中文，覆盖 GFM、高亮、TOC、stats 四条路径。
const sampleMD = "# 标题一\n\n" +
	"这是一段中文测试内容，用于验证字数统计。\n\n" +
	"## 二级标题\n\n" +
	"- 列表项一\n" +
	"- 列表项二\n\n" +
	"```go\n" +
	"func main() { println(\"hi\") }\n" +
	"```\n"

func TestRender(t *testing.T) {
	r := Render([]byte(sampleMD))
	if r == nil {
		t.Fatal("Render returned nil")
	}
	if r.HTML == "" {
		t.Fatal("HTML empty")
	}
	if !strings.Contains(r.HTML, "<h1") {
		t.Errorf("HTML missing heading, got: %s", r.HTML)
	}
	if !strings.Contains(r.HTML, "<pre") {
		// 代码块经 chroma 高亮后输出 <pre>
		t.Errorf("HTML missing highlighted code block, got: %s", r.HTML)
	}
	if !strings.Contains(r.TOC, "<ol") {
		t.Errorf("TOC missing <ol>, got: %s", r.TOC)
	}
	if r.Stats == nil {
		t.Fatal("Stats nil")
	}
	if r.Stats.Words == 0 {
		t.Errorf("Stats.Words == 0")
	}
}

func TestRenderEmpty(t *testing.T) {
	r := Render(nil)
	if r.HTML != "" || r.TOC != "" || r.Stats != nil {
		t.Errorf("empty input should yield zero Result, got %+v", r)
	}
}

func TestStatsCompatibility(t *testing.T) {
	r := Render([]byte("# Hello 世界 123\n\nHello"))
	if r.Stats == nil {
		t.Fatal("Stats nil")
	}
	if r.Stats.Words != 4 || r.Stats.Chars != 17 || r.Stats.NoSpaces != 15 {
		t.Errorf("unexpected stats: %+v", r.Stats)
	}
	if r.Stats.Frequency["hello"] != 2 || r.Stats.Frequency["世界"] != 1 || r.Stats.Frequency["123"] != 1 {
		t.Errorf("unexpected frequency: %#v", r.Stats.Frequency)
	}
}

func TestToHTML(t *testing.T) {
	if got := ToHTML([]byte(sampleMD)); !strings.Contains(got, "<h1") {
		t.Errorf("ToHTML missing heading, got: %s", got)
	}
	if got := ToHTML(nil); got != "" {
		t.Errorf("ToHTML(nil) should be empty, got: %s", got)
	}
}

func TestCountWords(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int
	}{
		{"空内容", "", 0},
		{"中文逐字计数", "你好世界", 4},
		{"英文按字符计数", "hello world", 10},
		{"中英混排", "Go 是一门编译型语言", 10},
		{"标题正文与标点", "# Hello 世界 123\n\nHello", 15},
		{"标点与空白不计入", "春眠不觉晓，处处闻啼鸟。", 10},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := CountWords([]byte(c.in)); got != c.want {
				t.Errorf("CountWords(%q) = %d, want %d", c.in, got, c.want)
			}
		})
	}
}

func TestLazyImages(t *testing.T) {
	r := Render([]byte("![示例图](/images/a.png)\n\n正文段落。\n"))
	if !strings.Contains(r.HTML, `loading="lazy"`) {
		t.Errorf("img missing loading=lazy, got: %s", r.HTML)
	}
	if !strings.Contains(r.HTML, `decoding="async"`) {
		t.Errorf("img missing decoding=async, got: %s", r.HTML)
	}
	if !strings.Contains(r.HTML, `alt="示例图"`) {
		t.Errorf("img alt lost, got: %s", r.HTML)
	}
}
