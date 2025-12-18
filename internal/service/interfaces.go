package service

import (
	"gitee.com/jieepre/go-site/internal/model"
)

// PostServiceInterface 文章服务接口
type PostServiceInterface interface {
	// GetPostBySlug 根据 slug 获取文章
	GetPostBySlug(slug string) (*model.Post, error)
	// GetPostList 获取文章列表
	GetPostList(pageNum, pageSize int) ([]model.LatestPosts, int64, error)
	// CreatePost 创建文章
	CreatePost(post *model.Post) error
	// UpdatePost 更新文章
	UpdatePost(post *model.Post) error
	// DeletePost 删除文章（软删除）
	DeletePost(id uint64) error
}

// CategoryServiceInterface 分类服务接口
type CategoryServiceInterface interface {
	// GetAllCategories 获取所有分类
	GetAllCategories() ([]model.Category, error)
	// GetCategoryByID 根据 ID 获取分类
	GetCategoryByID(id uint32) (*model.Category, error)
	// CreateCategory 创建分类
	CreateCategory(category *model.Category) error
	// UpdateCategory 更新分类
	UpdateCategory(category *model.Category) error
	// DeleteCategory 删除分类
	DeleteCategory(id uint32) error
}

// TagServiceInterface 标签服务接口
type TagServiceInterface interface {
	// GetAllTags 获取所有标签
	GetAllTags() ([]model.Tag, error)
	// GetTagByID 根据 ID 获取标签
	GetTagByID(id uint32) (*model.Tag, error)
	// CreateTag 创建标签
	CreateTag(tag *model.Tag) error
	// DeleteTag 删除标签
	DeleteTag(id uint32) error
}

// LinkServiceInterface 友情链接服务接口
type LinkServiceInterface interface {
	// GetAllLinks 获取所有链接
	GetAllLinks() ([]model.FriendLink, error)
	// CreateLink 创建链接
	CreateLink(link *model.FriendLink) error
	// UpdateLink 更新链接
	UpdateLink(link *model.FriendLink) error
	// DeleteLink 删除链接
	DeleteLink(id uint32) error
}

// DashboardServiceInterface 仪表盘服务接口
type DashboardServiceInterface interface {
	// GetDashboardStats 获取仪表盘统计
	GetDashboardStats() (*model.DashboardData, error)
}
