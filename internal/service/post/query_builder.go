package post

import (
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/querybuilder"
	"gorm.io/gorm"
)

// QueryBuilder 文章查询构建器
type QueryBuilder struct {
	*querybuilder.Builder
}

// NewQueryBuilder 创建文章查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{Builder: querybuilder.New(db, model.TPostsTable)}
}

func (qb *QueryBuilder) Select(fields string) *QueryBuilder {
	qb.Builder.Select(fields)
	return qb
}
func (qb *QueryBuilder) WithPublished() *QueryBuilder {
	qb.Builder.Scope(PublishedScope())
	return qb
}
func (qb *QueryBuilder) WithNotDeleted() *QueryBuilder {
	qb.Builder.Scope(NotDeletedScope())
	return qb
}
func (qb *QueryBuilder) WithCategory() *QueryBuilder {
	qb.Builder.Scope(WithCategoryScope())
	return qb
}
func (qb *QueryBuilder) WithTags() *QueryBuilder {
	qb.Builder.Scope(WithTagsScope())
	return qb
}
func (qb *QueryBuilder) WithPagination(pageNum, pageSize int) *QueryBuilder {
	qb.Builder.Scope(PaginationScope(pageNum, pageSize))
	return qb
}
func (qb *QueryBuilder) WithOrderByTime(desc bool) *QueryBuilder {
	qb.Builder.Scope(OrderByCreateTimeScope(desc))
	return qb
}
func (qb *QueryBuilder) WithOrderByTopAndTime() *QueryBuilder {
	qb.Builder.Scope(OrderByTopAndTimeScope())
	return qb
}
func (qb *QueryBuilder) ById(id int) *QueryBuilder {
	qb.Builder.Scope(PostByIdScope(id))
	return qb
}
func (qb *QueryBuilder) ByCategory(categoryName string) *QueryBuilder {
	qb.Builder.Scope(PostByCategoryScope(categoryName))
	return qb
}
func (qb *QueryBuilder) ByCategoryId(categoryID int64) *QueryBuilder {
	qb.Builder.Scope(PostByCategoryIdScope(categoryID))
	return qb
}
func (qb *QueryBuilder) ByTag(tagName string) *QueryBuilder {
	qb.Builder.Scope(PostByTagScope(tagName))
	return qb
}
func (qb *QueryBuilder) WithSearch(keyword string) *QueryBuilder {
	qb.Builder.Scope(SearchScope(keyword))
	return qb
}
func (qb *QueryBuilder) WithTitle(title string) *QueryBuilder {
	qb.Builder.Scope(PostTitleScope(title))
	return qb
}
func (qb *QueryBuilder) WithPublishedStatus(published *uint8) *QueryBuilder {
	qb.Builder.Scope(PostPublishedStatusScope(published))
	return qb
}
func (qb *QueryBuilder) WithCategoryId(categoryID *uint32) *QueryBuilder {
	qb.Builder.Scope(PostCategoryIdScope(categoryID))
	return qb
}

func (qb *QueryBuilder) Build() *gorm.DB          { return qb.Builder.Build() }
func (qb *QueryBuilder) Count(count *int64) error { return qb.Builder.Count(count) }
func (qb *QueryBuilder) Find(dest any) error      { return qb.Builder.Find(dest) }
func (qb *QueryBuilder) First(dest any) error     { return qb.Builder.First(dest) }
func (qb *QueryBuilder) Scan(dest any) error      { return qb.Builder.Scan(dest) }
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.Builder.Limit(limit)
	return qb
}
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.Builder.Offset(offset)
	return qb
}
func (qb *QueryBuilder) Order(order string) *QueryBuilder {
	qb.Builder.Order(order)
	return qb
}
