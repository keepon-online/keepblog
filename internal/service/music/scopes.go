package music

import (
	"gorm.io/gorm"
)

// ActiveScope 只查询启用状态的音乐
func ActiveScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("state = ?", 1)
	}
}

// MusicByIdScope 根据 ID 查询音乐
func MusicByIdScope(id uint32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// OrderBySortScope 按排序字段排序
func OrderBySortScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC, id DESC")
	}
}
