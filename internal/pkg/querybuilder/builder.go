// Package querybuilder 提供各领域查询构建器共享的 GORM 基础能力。
// 领域包只负责表名和业务 scope，避免重复维护相同的 CRUD 链式实现。
package querybuilder

import "gorm.io/gorm"

// Builder 是可组合的 GORM 查询基础构建器。
type Builder struct {
	DB     *gorm.DB
	Scopes []func(*gorm.DB) *gorm.DB
}

// New 创建以 db 为基础的构建器。
func New(db *gorm.DB, table string) *Builder {
	return &Builder{DB: db.Table(table), Scopes: make([]func(*gorm.DB) *gorm.DB, 0)}
}

// Select 设置查询字段。
func (b *Builder) Select(fields string) *Builder {
	b.DB = b.DB.Select(fields)
	return b
}

// Scope 追加 GORM scope。
func (b *Builder) Scope(scope func(*gorm.DB) *gorm.DB) *Builder {
	b.Scopes = append(b.Scopes, scope)
	return b
}

// Build 应用全部 scope。
func (b *Builder) Build() *gorm.DB {
	return b.DB.Scopes(b.Scopes...)
}

func (b *Builder) Count(count *int64) error {
	return b.Build().Count(count).Error
}

func (b *Builder) Find(dest any) error {
	return b.Build().Find(dest).Error
}

func (b *Builder) First(dest any) error {
	return b.Build().First(dest).Error
}

func (b *Builder) Scan(dest any) error {
	return b.Build().Scan(dest).Error
}

func (b *Builder) Limit(limit int) *Builder {
	b.DB = b.DB.Limit(limit)
	return b
}

func (b *Builder) Offset(offset int) *Builder {
	b.DB = b.DB.Offset(offset)
	return b
}

func (b *Builder) Order(order string) *Builder {
	b.DB = b.DB.Order(order)
	return b
}
