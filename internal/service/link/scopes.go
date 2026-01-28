package link

import (
	"gorm.io/gorm"
)

// ActiveScope 只查询启用状态的友链
func ActiveScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("state = ?", 1)
	}
}

// LinkByIdScope 根据 ID 查询友链
func LinkByIdScope(id uint32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// OrderByCreateTimeScope 按创建时间排序
func OrderByCreateTimeScope(desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order("create_time DESC")
		}
		return db.Order("create_time ASC")
	}
}
