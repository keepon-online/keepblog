package category

import (
	"gorm.io/gorm"

	"gitee.com/jieepre/keepblog/internal/pkg/querybuilder"
)

// ActiveScope 只查询启用状态的分类
func ActiveScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category.state = ?", 1)
	}
}

// CategoryByIdScope 根据 ID 查询分类
func CategoryByIdScope(id uint32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category.category_id = ?", id)
	}
}

// WithPostCountScope 关联文章数量统计
func WithPostCountScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Select("category.category_name, COUNT(post.post_id) as total").
			Joins("LEFT JOIN post ON category.category_id = post.category_id").
			Where(querybuilder.VisibleWhere("post"), querybuilder.VisibleNow()).
			Where("post.is_deleted = ?", 0).
			Group("category.category_name")
	}
}
