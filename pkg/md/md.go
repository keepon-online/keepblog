package md

import (
	"bytes"
	"strings"
	"sync"

	formathtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/gookit/slog"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/toc"
)

// 高亮样式，集中一处便于统一调整。
const highlightStyle = "monokai"

var (
	// markdown 负责解析 Markdown 与渲染正文 HTML：带 GFM、代码高亮、自动标题 ID。
	// 正文不再内嵌目录，侧边栏目录由 Render 通过 toc.Inspect 单独产出。
	// goldmark.Markdown 本身是 goroutine 安全的，构造一次后即可并发复用。
	markdown goldmark.Markdown
	once     sync.Once
)

// initEngines 懒初始化 goldmark 实例。
func initEngines() {
	once.Do(func() {
		markdown = goldmark.New(
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
				parser.WithAttribute(),
				parser.WithHeadingAttribute(),
			),
			goldmark.WithRendererOptions(html.WithHardWraps()),
			goldmark.WithExtensions(
				extension.GFM,
				extension.Table,
				extension.Linkify,
				highlighting.NewHighlighting(
					highlighting.WithStyle(highlightStyle),
					highlighting.WithFormatOptions(
						formathtml.LinkableLineNumbers(true, ""),
						formathtml.WithLineNumbers(true),
					),
				),
			),
		)
	})
}

// Result 是一次 Markdown 解析产出的集合，避免对同一份内容多次解析。
type Result struct {
	HTML  string // 渲染后的正文 HTML
	TOC   string // 侧边栏目录（<ol>...</ol>），无标题时为空
	Stats *Info  // 字数/阅读时长统计，content 为空时为 nil
}

// Render 一次解析 content，同时产出正文 HTML、侧边栏目录和字数统计。
// 适用于文章页这类同时需要三者的场景。content 为空时返回零值 Result。
func Render(content []byte) *Result {
	r := &Result{}
	if len(content) == 0 {
		return r
	}
	initEngines()

	// 正文 HTML。
	var buf bytes.Buffer
	if err := markdown.Convert(content, &buf); err != nil {
		slog.Errorf("md to html fail: %s", err.Error())
		return r
	}
	r.HTML = buf.String()

	// 目录与统计共享同一次解析。
	doc := markdown.Parser().Parse(text.NewReader(content))
	r.Stats = NewStats(doc, content)

	if tree, err := toc.Inspect(doc, content); err != nil {
		slog.Errorf("md toc inspect fail: %s", err.Error())
	} else if treeList := toc.RenderList(tree); treeList != nil {
		var bufToc bytes.Buffer
		if err := markdown.Renderer().Render(&bufToc, content, treeList); err != nil {
			slog.Errorf("md toc render fail: %s", err.Error())
		} else {
			// 侧边栏模板（card-toc.html）使用 <ol>，与 goldmark 默认的 <ul> 对齐。
			r.TOC = strings.ReplaceAll(bufToc.String(), "ul", "ol")
		}
	}
	return r
}

// ToHTML 仅渲染正文 HTML，不计算目录与统计。适用于只需正文的页面（如关于页）。
func ToHTML(content []byte) string {
	if len(content) == 0 {
		return ""
	}
	initEngines()
	var buf bytes.Buffer
	if err := markdown.Convert(content, &buf); err != nil {
		slog.Errorf("md to html fail: %s", err.Error())
		return ""
	}
	return buf.String()
}

// Goldmark2html 渲染正文 HTML。
//
// Deprecated: 改用 ToHTML（仅 HTML）或 Render（HTML+TOC+Stats）。
func Goldmark2html(content []byte) []byte {
	return []byte(ToHTML(content))
}

// Goldmark2htmlToc 渲染侧边栏目录（<ol>）。
//
// Deprecated: 改用 Render，可一次解析同时拿到 HTML/TOC/Stats。
func Goldmark2htmlToc(content []byte) string {
	if len(content) == 0 {
		return ""
	}
	return Render(content).TOC
}

// Goldmarkstats 返回字数/阅读时长统计。
//
// Deprecated: 改用 Render，可一次解析同时拿到 HTML/TOC/Stats。
func Goldmarkstats(content []byte) *Info {
	if len(content) == 0 {
		return nil
	}
	return Render(content).Stats
}
