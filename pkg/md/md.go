package md

import (
	"bytes"
	"strings"
	"sync"

	formathtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/gookit/slog"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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

// lazyImages 为正文图片补 loading="lazy" decoding="async"。
// 正文图绝大多数在首屏之外，懒加载可减少初始请求、改善 LCP；
// 若首图恰在首屏，浏览器对首屏内的 lazy 图仍会立即加载，无副作用。
type lazyImages struct{}

func (lazyImages) Transform(node *gast.Document, _ text.Reader, _ parser.Context) {
	_ = gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		// Walk 对每个节点进出各回调一次，只在进入时处理避免重复设置。
		if entering {
			if img, ok := n.(*gast.Image); ok {
				img.SetAttributeString("loading", []byte("lazy"))
				img.SetAttributeString("decoding", []byte("async"))
			}
		}
		return gast.WalkContinue, nil
	})
}

// initEngines 懒初始化 goldmark 实例。
func initEngines() {
	once.Do(func() {
		markdown = goldmark.New(
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
				parser.WithAttribute(),
				parser.WithHeadingAttribute(),
				parser.WithASTTransformers(
					util.Prioritized(lazyImages{}, 100),
				),
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

// CountWords 统计 Markdown 正文的字数：只解析 AST 不渲染 HTML，
// 适用于入库时计算 word_count。口径与 Info.NoSpaces 一致——仅统计
// 字母与数字（CJK 文字逐字计数，空白与标点不计入）。
func CountWords(content []byte) int {
	if len(content) == 0 {
		return 0
	}
	initEngines()
	doc := markdown.Parser().Parse(text.NewReader(content))
	return NewStats(doc, content).NoSpaces
}

// Excerpt 从 Markdown 提取纯文本摘要：跳过代码块与行内代码，
// 收集正文文本后按 rune 截取 limit 字（中文按字数截断）。
// 用作文章 summary 为空时的 meta description / RSS description 兜底。
func Excerpt(content []byte, limit int) string {
	if len(content) == 0 || limit <= 0 {
		return ""
	}
	initEngines()
	doc := markdown.Parser().Parse(text.NewReader(content))

	var sb strings.Builder
	_ = gast.Walk(doc, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		switch n.Kind() {
		case gast.KindCodeBlock, gast.KindFencedCodeBlock, gast.KindCodeSpan:
			return gast.WalkSkipChildren, nil
		case gast.KindText, gast.KindString:
			sb.Write(n.Text(content))
		}
		return gast.WalkContinue, nil
	})

	runes := []rune(strings.TrimSpace(sb.String()))
	if len(runes) > limit {
		runes = runes[:limit]
	}
	return string(runes)
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
