package post

import (
	"testing"

	"gitee.com/jieepre/go-site/global"
)

// TestQueryBuilder 测试查询构建器。
// 注意：QueryBuilder 在构造时即持有 global.GORM（设计上依赖运行时 DB 连接），
// 单元测试若未初始化 GORM 会 panic，故无 DB 时跳过。
func TestQueryBuilder(t *testing.T) {
	if global.GORM == nil {
		t.Skip("跳过：global.GORM 未初始化（需真实数据库连接）")
	}

	// 测试基本构建
	qb := NewQueryBuilder()
	if qb == nil {
		t.Error("NewQueryBuilder() 返回 nil")
	}

	// 测试链式调用
	qb = NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		WithCategory().
		WithTags().
		WithPagination(1, 10)

	if qb == nil {
		t.Error("链式调用失败")
	}

	// 测试 Build
	db := qb.Build()
	if db == nil {
		t.Error("Build() 返回 nil")
	}
}

// TestScopes 测试 Scopes
func TestScopes(t *testing.T) {
	// 测试 PublishedScope
	scope := PublishedScope()
	if scope == nil {
		t.Error("PublishedScope() 返回 nil")
	}

	// 测试 NotDeletedScope
	scope = NotDeletedScope()
	if scope == nil {
		t.Error("NotDeletedScope() 返回 nil")
	}

	// 测试 PaginationScope
	scope = PaginationScope(1, 10)
	if scope == nil {
		t.Error("PaginationScope() 返回 nil")
	}
}

// TestServiceMethods 测试服务方法存在性
func TestServiceMethods(t *testing.T) {
	// NewPostService 总是返回非 nil 的 &Service{}，这里取方法引用做编译期检查。
	service := NewPostService()

	// 验证方法存在（编译时检查）
	_ = service.GetPost
	_ = service.GetPostDetail
	_ = service.GetLatestPosts
	_ = service.GetCoverPosts
	_ = service.GetPostsByCategory
	_ = service.GetPostsByTag
	_ = service.Search
	_ = service.SearchPaged
	_ = service.SavePost
	_ = service.UpdatePost
	_ = service.DeletePost
}
