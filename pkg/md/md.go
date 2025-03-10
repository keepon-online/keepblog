package md

import (
	"bytes"
	"fmt"
	"github.com/88250/lute"
	formathtml "github.com/alecthomas/chroma/formatters/html"
	stats "github.com/mdigger/goldmark-stats"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/toc"
	"strings"
)

func Goldmark2html(content []byte) []byte {

	markdown := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithHeadingAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
		// 支持 GFM
		goldmark.WithExtensions(extension.GFM, extension.Table, extension.Linkify),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					formathtml.LinkableLineNumbers(true, ""),
					formathtml.WithLineNumbers(true),
				),
			),
		),
	)
	markdown.Parser().AddOptions(
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
		parser.WithHeadingAttribute(),
		parser.WithASTTransformers(
			util.Prioritized(&toc.Transformer{
				Title: "Contents",
			}, 100),
		),
	)
	var buf bytes.Buffer

	if err := markdown.Convert(content, &buf); err != nil {

		fmt.Println("md to html fail," + err.Error())
		return nil
	}

	return buf.Bytes()
}

func Goldmark2htmlToc(content []byte) string {
	if content == nil {
		return ""
	}
	markdown := goldmark.New()
	markdown.Parser().AddOptions(parser.WithAutoHeadingID(), parser.WithAttribute(), parser.WithHeadingAttribute())
	doc := markdown.Parser().Parse(text.NewReader(content))
	tree, err := toc.Inspect(doc, content)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	treeList := toc.RenderList(tree)
	var bufToc bytes.Buffer
	err = markdown.Renderer().Render(&bufToc, content, treeList)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return strings.ReplaceAll(bufToc.String(), "ul", "ol")
}

func Goldmarkstats(content []byte) *stats.Info {
	if content == nil {
		return nil
	}
	doc := goldmark.DefaultParser().Parse(text.NewReader(content))
	return stats.New(doc, content)
}

func LuneTet(content string) string {
	luteEngine := lute.New() // 默认已经启用 GFM 支持以及中文语境优化
	luteEngine.SetCodeSyntaxHighlightStyleName("xcode-dark")
	luteEngine.SetCodeSyntaxHighlightLineNum(true)
	luteEngine.SetCodeSyntaxHighlight(true)
	luteEngine.SetCodeSyntaxHighlightInlineStyle(true)
	//luteEngine.SetCodeSyntaxHighlightDetectLang(true)
	luteEngine.SetToC(true)
	return luteEngine.MarkdownStr("demo", content)
}
