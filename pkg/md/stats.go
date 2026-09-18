package md

import (
	"bytes"
	"fmt"
	"sort"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// Info stats from markdown text.
type Info struct {
	Words     int
	Chars     int
	NoSpaces  int
	Frequency map[string]int
	mostUsed  []WordTuple
}

func (info Info) Duration(speed int) time.Duration {
	return (time.Duration(info.Chars) * time.Minute / time.Duration(speed)).Round(time.Second)
}

func (info Info) Pages() int {
	return (info.Chars + 1800 - 1) / 1800
}

func (info Info) AuthorPages() float32 {
	return float32(info.Chars) / 40000
}

func (info Info) Unique() int {
	return len(info.Frequency)
}

func (info *Info) MostUsed() []WordTuple {
	if info.mostUsed == nil {
		info.mostUsed = make([]WordTuple, 0, info.Unique())
		for word, count := range info.Frequency {
			info.mostUsed = append(info.mostUsed, WordTuple{Word: word, Count: count})
		}
		sort.Slice(info.mostUsed, func(i, j int) bool {
			if info.mostUsed[i].Count == info.mostUsed[j].Count {
				return info.mostUsed[i].Word < info.mostUsed[j].Word
			}
			return info.mostUsed[i].Count > info.mostUsed[j].Count
		})
	}
	return info.mostUsed
}

type WordTuple struct {
	Word  string
	Count int
}

func (t WordTuple) String() string {
	return fmt.Sprintf("%v: %v", t.Word, t.Count)
}

func NewStats(node ast.Node, source []byte) *Info {
	var words, chars, noSpaces int
	frequency := make(map[string]int, 10000)
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && (node.Kind() == ast.KindText || node.Kind() == ast.KindString) {
			text := node.Text(source)
			chars += utf8.RuneCount(text)
			for _, word := range bytes.FieldsFunc(text, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			}) {
				words++
				noSpaces += utf8.RuneCount(word)
				frequency[util.BytesToReadOnlyString(bytes.ToLower(word))]++
			}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return &Info{Words: words, Chars: chars, NoSpaces: noSpaces, Frequency: frequency}
}
