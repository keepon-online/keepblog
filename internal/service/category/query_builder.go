package category

import (
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/querybuilder"
	"gorm.io/gorm"
)

// QueryBuilder 分类查询构建器
type QueryBuilder struct {
	*querybuilder.Builder
}

// NewQueryBuilder 创建分类查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{Builder: querybuilder.New(db, model.TCategoryTable)}
}

func (qb *QueryBuilder) Select(fields string) *QueryBuilder {
	qb.Builder.Select(fields)
	return qb
}

func (qb *QueryBuilder) WithActive() *QueryBuilder {
	qb.Builder.Scope(ActiveScope())
	return qb
}

func (qb *QueryBuilder) ById(id uint32) *QueryBuilder {
	qb.Builder.Scope(CategoryByIdScope(id))
	return qb
}

func (qb *QueryBuilder) WithPostCount() *QueryBuilder {
	qb.Builder.Scope(WithPostCountScope())
	return qb
}

func (qb *QueryBuilder) Build() *gorm.DB          { return qb.Builder.Build() }
func (qb *QueryBuilder) Count(count *int64) error { return qb.Builder.Count(count) }
func (qb *QueryBuilder) Find(dest any) error      { return qb.Builder.Find(dest) }
func (qb *QueryBuilder) First(dest any) error     { return qb.Builder.First(dest) }
