package post

import (
	"gorm.io/gorm"
)

// 常量定义
const (
	DefaultPageSize = 10
)

// PublishedScope 只查询已发布的文章
func PublishedScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("post.is_published = ?", 1)
	}
}

// NotDeletedScope 只查询未删除的文章
func NotDeletedScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("post.is_deleted = ?", 0)
	}
}

// WithCategoryScope 关联分类表
func WithCategoryScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("LEFT JOIN category ON post.category_id = category.category_id")
	}
}

// WithTagsScope 关联标签表
func WithTagsScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("LEFT JOIN post_tag pt ON post.post_id = pt.post_id").
			Joins("LEFT JOIN tag t ON pt.tag_id = t.tag_id")
	}
}

// PaginationScope 分页
func PaginationScope(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageSize <= 0 {
			pageSize = DefaultPageSize
		}
		if pageNum <= 0 {
			pageNum = 1
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// OrderByCreateTimeScope 按创建时间排序
func OrderByCreateTimeScope(desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order("post.create_time DESC")
		}
		return db.Order("post.create_time ASC")
	}
}

// OrderByTopAndTimeScope 按置顶和时间排序
func OrderByTopAndTimeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("post.top DESC, post.create_time DESC")
	}
}

// PostByIdScope 根据 ID 查询
func PostByIdScope(id int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("post.post_id = ?", id)
	}
}

// PostByCategoryScope 根据分类查询
func PostByCategoryScope(categoryName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if categoryName == "" {
			return db
		}
		return db.Where("category.category_name = ?", categoryName)
	}
}

// PostByCategoryIdScope 根据分类 ID 查询
func PostByCategoryIdScope(categoryId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if categoryId <= 0 {
			return db
		}
		return db.Where("post.category_id = ?", categoryId)
	}
}

// PostByTagScope 根据标签查询
func PostByTagScope(tagName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if tagName == "" {
			return db
		}
		return db.Where("t.tag_name = ?", tagName)
	}
}

// SearchScope 搜索
func SearchScope(keyword string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" {
			return db
		}
		likeKeyword := "%" + keyword + "%"
		return db.Where(
			"post.title LIKE ? OR post.summary LIKE ? OR post.post_content LIKE ?",
			likeKeyword, likeKeyword, likeKeyword,
		)
	}
}

// PostTitleScope 根据标题模糊查询
func PostTitleScope(title string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if title == "" {
			return db
		}
		return db.Where("post.title LIKE ?", "%"+title+"%")
	}
}

// PostPublishedStatusScope 根据发布状态查询
func PostPublishedStatusScope(published *uint8) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if published == nil {
			return db
		}
		return db.Where("post.is_published = ?", *published)
	}
}

// PostCategoryIdScope 根据分类 ID 查询（用于后台筛选）
func PostCategoryIdScope(categoryId *uint32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if categoryId == nil {
			return db
		}
		return db.Where("post.category_id = ?", *categoryId)
	}
}
