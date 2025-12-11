package core

import (
	"gitee.com/jieepre/go-site/global"
	"github.com/gookit/slog"
)

// InitIndexes 初始化数据库索引
func InitIndexes() {
	db := global.GORM

	// 为经常查询的字段添加索引
	indexes := []string{
		// 文章表索引
		"CREATE INDEX IF NOT EXISTS idx_post_published ON post (is_published);",
		"CREATE INDEX IF NOT EXISTS idx_post_deleted ON post (is_deleted);",
		"CREATE INDEX IF NOT EXISTS idx_post_category ON post (category_id);",
		"CREATE INDEX IF NOT EXISTS idx_post_status ON post (status);",
		"CREATE INDEX IF NOT EXISTS idx_post_create_time ON post (create_time);",
		"CREATE INDEX IF NOT EXISTS idx_post_pub_time ON post (pub_time);",
		"CREATE INDEX IF NOT EXISTS idx_post_slug ON post (post_slug);",

		// 复合索引 - 文章列表查询优化
		"CREATE INDEX IF NOT EXISTS idx_post_list ON post (is_published, is_deleted, create_time DESC);",

		// 分类表索引
		"CREATE INDEX IF NOT EXISTS idx_category_state ON category (state);",

		// 标签表索引
		"CREATE INDEX IF NOT EXISTS idx_tag_name ON tag (tag_name);",

		// 文章标签关联表索引
		"CREATE INDEX IF NOT EXISTS idx_post_tag_post ON post_tag (post_id);",
		"CREATE INDEX IF NOT EXISTS idx_post_tag_tag ON post_tag (tag_id);",

		// 评论表索引
		"CREATE INDEX IF NOT EXISTS idx_comment_post ON comment (post_id);",
		"CREATE INDEX IF NOT EXISTS idx_comment_approved ON comment (is_approved);",
		"CREATE INDEX IF NOT EXISTS idx_comment_create_time ON comment (create_time);",

		// 访问日志索引
		"CREATE INDEX IF NOT EXISTS idx_access_log_ip ON system_access_log (ip);",
		"CREATE INDEX IF NOT EXISTS idx_access_log_url ON system_access_log (url);",
		"CREATE INDEX IF NOT EXISTS idx_access_log_create_time ON system_access_log (create_at);",
		"CREATE INDEX IF NOT EXISTS idx_access_log_status ON system_access_log (status);",

		// 登录日志索引
		"CREATE INDEX IF NOT EXISTS idx_login_log_create_time ON system_login_log (create_at);",
		"CREATE INDEX IF NOT EXISTS idx_login_log_success ON system_login_log (success);",

		// 友情链接索引
		"CREATE INDEX IF NOT EXISTS idx_friend_link_state ON friend_link (state);",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			slog.Errorf("Failed to create index: %s, error: %v", indexSQL, err)
		}
	}

	slog.Info("Database indexes initialized successfully")
}
