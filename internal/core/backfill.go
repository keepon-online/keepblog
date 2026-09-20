package core

import (
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/pkg/md"
	"github.com/gookit/slog"
	"gorm.io/gorm"
)

// BackfillPostWordCount 回填存量文章的 word_count。
// word_count 列自建库起存在但从未被写入，2026-09 起保存/更新文章时由服务端
// 计算入库；此函数在启动时一次性补齐旧数据。幂等：只处理 word_count 为 0
// 且正文非空的文章，正文为空的重复计算代价可忽略。
// 使用 UpdateColumn 而非 Update，避免触碰 last_modified_time（它驱动
// 侧边栏"最后更新时间"的显示）。
func BackfillPostWordCount(db *gorm.DB) {
	if db == nil {
		return
	}

	var posts []model.Post
	if err := db.Table(model.TPostsTable).
		Select("post_id, post_content").
		Where("word_count = 0 AND post_content <> ''").
		Find(&posts).Error; err != nil {
		slog.Errorf("回填字数：查询待回填文章失败: %v", err)
		return
	}

	updated := 0
	for _, p := range posts {
		words := md.CountWords([]byte(p.PostContent))
		if words == 0 {
			continue
		}
		if err := db.Table(model.TPostsTable).
			Where("post_id = ?", p.PostId).
			UpdateColumn("word_count", words).Error; err != nil {
			slog.Errorf("回填字数：文章 %d 更新失败: %v", p.PostId, err)
			continue
		}
		updated++
	}
	if updated > 0 {
		slog.Infof("回填字数：已补齐 %d 篇文章的 word_count", updated)
	}
}
