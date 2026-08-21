package link

import (
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/querybuilder"
	"gorm.io/gorm"
)

// QueryBuilder 友链查询构建器
type QueryBuilder struct {
	*querybuilder.Builder
}

// NewQueryBuilder 创建友链查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{Builder: querybuilder.New(db, model.TFriendLinkTable)}
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
	qb.Builder.Scope(LinkByIdScope(id))
	return qb
}

func (qb *QueryBuilder) WithOrderByTime(desc bool) *QueryBuilder {
	qb.Builder.Scope(OrderByCreateTimeScope(desc))
	return qb
}

func (qb *QueryBuilder) Build() *gorm.DB          { return qb.Builder.Build() }
func (qb *QueryBuilder) Count(count *int64) error { return qb.Builder.Count(count) }
func (qb *QueryBuilder) Find(dest any) error      { return qb.Builder.Find(dest) }
func (qb *QueryBuilder) First(dest any) error     { return qb.Builder.First(dest) }
