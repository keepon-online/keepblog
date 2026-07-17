package core

import (
	"html/template"
	"strings"
	"time"
)

func TemplateFunc() template.FuncMap {
	return template.FuncMap{
		"unescaped": func(str string) template.HTML {
			return template.HTML(str)
		},
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"css": func(s string) template.CSS {
			return template.CSS(s)
		},
		"js": func(s string) template.JS {
			return template.JS(s)
		},
		"url": func(s string) template.URL {
			return template.URL(s)
		},
		"time": func(t uint64) string {
			return time.Unix(int64(t), 0).Format(time.DateOnly)
		},
		"archiveTime": func(t string) string {
			// 兼容 "2026-07"（归档列表页 map key）与 "2026/07"（侧边栏）两种分隔符。
			parse, err := time.ParseInLocation("2006/01", strings.ReplaceAll(t, "-", "/"), time.Local)
			if err != nil {
				return t
			}
			return parse.Format("2006年01月")
		},
		"even": func(i int) bool {
			return i%2 == 0
		},
		"BaiduStat": func(t string) bool {
			return t != ""
		},
		"BaiduSite": func(t string) bool {
			return t != ""
		},
		"divf": func(a, b int64) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"gt": func(a, b int64) bool {
			return a > b
		},
	}
}
