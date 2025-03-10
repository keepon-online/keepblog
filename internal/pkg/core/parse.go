package core

import (
	"html/template"
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
			return time.Unix(int64(t), 0).Format("2006-01-02")
		},
		"archiveTime": func(t string) string {
			parse, _ := time.ParseInLocation("2006/01", t, time.Local)
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
	}
}
