package querybuilder

import (
	"fmt"
	"time"
)

// VisibleWhere 返回文章"已发布且发布时间已到"的 SQL 条件，alias 为 post 表
// 别名（"p"、"post"，空串表示无别名）。占位符 ? 需以 time.Now().Unix() 为参数传入。
//
// 定时发布语义：is_published=1 且 pub_time 在未来的文章到点前对前台不可见。
// 该条件被 post/sidebar/tag/category 与搜索引擎推送等多处共用，任何一处漏改
// 都会导致未到点文章从侧栏计数、标签页、sitemap 或推送中提前泄露。
func VisibleWhere(alias string) string {
	if alias != "" {
		alias += "."
	}
	return fmt.Sprintf("%sis_published = 1 AND (%spub_time IS NULL OR %spub_time <= ?)", alias, alias, alias)
}

// VisibleNow 返回 VisibleWhere 的当前时间参数，配合 gorm Where(cond, VisibleNow()) 使用。
func VisibleNow() int64 {
	return time.Now().Unix()
}
