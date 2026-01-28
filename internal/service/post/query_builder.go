package post

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gorm.io/gorm"
)

// QueryBuilder 查询构建器
type QueryBuilder struct {
	db     *gorm.DB
	scopes []func(*gorm.DB) *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		db:     global.GORM.Table(model.TPostsTable),
		scopes: make([]func(*gorm.DB) *gorm.DB, 0),
	}
}

// Select 设置查询字段
func (qb *QueryBuilder) Select(fields string) *QueryBuilder {
	qb.db = qb.db.Select(fields)
	return qb
}

// WithPublished 只查询已发布
func (qb *QueryBuilder) WithPublished() *QueryBuilder {
	qb.scopes = append(qb.scopes, PublishedScope())
	return qb
}

// WithNotDeleted 只查询未删除
func (qb *QueryBuilder) WithNotDeleted() *QueryBuilder {
	qb.scopes = append(qb.scopes, NotDeletedScope())
	return qb
}

// WithCategory 关联分类
func (qb *QueryBuilder) WithCategory() *QueryBuilder {
	qb.scopes = append(qb.scopes, WithCategoryScope())
	return qb
}

// WithTags 关联标签
func (qb *QueryBuilder) WithTags() *QueryBuilder {
	qb.scopes = append(qb.scopes, WithTagsScope())
	return qb
}

// WithPagination 分页
func (qb *QueryBuilder) WithPagination(pageNum, pageSize int) *QueryBuilder {
	qb.scopes = append(qb.scopes, PaginationScope(pageNum, pageSize))
	return qb
}

// WithOrderByTime 按时间排序
func (qb *QueryBuilder) WithOrderByTime(desc bool) *QueryBuilder {
	qb.scopes = append(qb.scopes, OrderByCreateTimeScope(desc))
	return qb
}

// WithOrderByTopAndTime 按置顶和时间排序
func (qb *QueryBuilder) WithOrderByTopAndTime() *QueryBuilder {
	qb.scopes = append(qb.scopes, OrderByTopAndTimeScope())
	return qb
}

// ById 根据 ID 查询
func (qb *QueryBuilder) ById(id int) *QueryBuilder {
	qb.scopes = append(qb.scopes, PostByIdScope(id))
	return qb
}

// ByCategory 根据分类查询
func (qb *QueryBuilder) ByCategory(categoryName string) *QueryBuilder {
	qb.scopes = append(qb.scopes, PostByCategoryScope(categoryName))
	return qb
}

// ByCategoryId 根据分类 ID 查询
func (qb *QueryBuilder) ByCategoryId(categoryId int64) *QueryBuilder {
	qb.scopes = append(qb.scopes, PostByCategoryIdScope(categoryId))
	return qb
}

// ByTag 根据标签查询
func (qb *QueryBuilder) ByTag(tagName string) *QueryBuilder {
	qb.scopes = append(qb.scopes, PostByTagScope(tagName))
	return qb
}

// WithSearch 搜索
func (qb *QueryBuilder) WithSearch(keyword string) *QueryBuilder {
	qb.scopes = append(qb.scopes, SearchScope(keyword))
	return qb
}

// WithTitle 根据标题查询
func (qb *QueryBuilder) WithTitle(title string) *QueryBuilder {
	qb.scopes = append(qb.scopes, PostTitleScope(title))
	return qb
}

// WithPublishedStatus 根据发布状态查询
func (qb *QueryBuilder) WithPublishedStatus(published *uint8) *QueryBuilder {
	qb.scopes = append(qb.scopes, PostPublishedStatusScope(published))
	return qb
}

// WithCategoryId 根据分类 ID 查询（用于后台筛选）
func (qb *QueryBuilder) WithCategoryId(categoryId *uint32) *QueryBuilder {
	qb.scopes = append(qb.scopes, PostCategoryIdScope(categoryId))
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

// Scan 扫描结果
func (qb *QueryBuilder) Scan(dest interface{}) error {
	return qb.Build().Scan(dest).Error
}

// Limit 限制数量
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.db = qb.db.Limit(limit)
	return qb
}

// Offset 偏移量
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.db = qb.db.Offset(offset)
	return qb
}

// Order 排序
func (qb *QueryBuilder) Order(order string) *QueryBuilder {
	qb.db = qb.db.Order(order)
	return qb
}
