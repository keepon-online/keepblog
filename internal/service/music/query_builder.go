package music

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gorm.io/gorm"
)

// QueryBuilder 音乐查询构建器
type QueryBuilder struct {
	db     *gorm.DB
	scopes []func(*gorm.DB) *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		db:     global.GORM.Table(model.TMusicTable),
		scopes: make([]func(*gorm.DB) *gorm.DB, 0),
	}
}

// Select 设置查询字段
func (qb *QueryBuilder) Select(fields string) *QueryBuilder {
	qb.db = qb.db.Select(fields)
	return qb
}

// WithActive 只查询启用状态
func (qb *QueryBuilder) WithActive() *QueryBuilder {
	qb.scopes = append(qb.scopes, ActiveScope())
	return qb
}

// ById 根据 ID 查询
func (qb *QueryBuilder) ById(id uint32) *QueryBuilder {
	qb.scopes = append(qb.scopes, MusicByIdScope(id))
	return qb
}

// WithOrderBySort 按排序字段排序
func (qb *QueryBuilder) WithOrderBySort() *QueryBuilder {
	qb.scopes = append(qb.scopes, OrderBySortScope())
	return qb
}

// Build 构建查询
func (qb *QueryBuilder) Build() *gorm.DB {
	return qb.db.Scopes(qb.scopes...)
}

// Count 统计数量
func (qb *QueryBuilder) Count(count *int64) error {
	return qb.Build().Count(count).Error
}

// Find 查询列表
func (qb *QueryBuilder) Find(dest interface{}) error {
	return qb.Build().Find(dest).Error
}

// First 查询单条
func (qb *QueryBuilder) First(dest interface{}) error {
	return qb.Build().First(dest).Error
}
