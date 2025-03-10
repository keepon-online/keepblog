package md

import (
	"bytes"
	"fmt"
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
	"os"
	"testing"
)

func TestGold(t *testing.T) {

	markdown := goldmark.New(
		goldmark.WithParserOptions(
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
		parser.WithASTTransformers(
			util.Prioritized(&toc.Transformer{
				Title: "Contents",
			}, 100),
		),
	)
	f, _ := os.Create("guide4.html")
	source, _ := os.ReadFile("README.md")

	//convertFunc := toc.Markdown(markdown)
	//
	//headers, _ := convertFunc(source, f)
	//
	//for _, header := range headers {
	//	fmt.Printf("%+v\n", header)
	//}

	// stats
	doc := goldmark.DefaultParser().Parse(text.NewReader(source))
	info := stats.New(doc, source)

	fmt.Printf("words: %d, unique: %d, chars: %d, reading time: %v\n",
		info.Words, info.Unique(), info.Chars, info.Duration(400))

	// Request that IDs are automatically assigned to headers.
	//markdown.Parser().AddOptions(parser.WithAutoHeadingID())
	// Alternatively, we can provide our own implementation of parser.IDs
	// and use,
	//
	//  pctx := parser.NewContext(parser.WithIDs(ids))
	//doc := parser.Parse(text.NewReader(src), parser.WithContext(pctx))

	// Inspect the parsed Markdown document to find headers and build a
	// tree for the table of contents.
	tree, err := toc.Inspect(doc, source)
	if err != nil {
		panic(err)
	}

	// Render the tree as-is into a Markdown list.
	treeList := toc.RenderList(tree)

	// Render the Markdown list into HTML.
	markdown.Renderer().Render(os.Stdout, source, treeList)
	var buf bytes.Buffer
	if err := markdown.Convert(source, f); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

}
