package post

import (
	"errors"
	"github.com/gookit/slog"

	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/request"
	"gitee.com/jieepre/go-site/pkg"
	"gitee.com/jieepre/go-site/pkg/hash"
	"gitee.com/jieepre/go-site/pkg/page"
)

// SavePost 保存文章
func (service Service) SavePost(content model.Post) (uint64, error) {
	if err := service.db.Table(model.TPostsTable).Create(&content).Error; err != nil {
		return content.PostId, errors.New("保存失败: " + err.Error())
	}
	hid, _ := hash.New().HashidsEncode([]int{int(content.PostId)})
	err := service.UpdatePostHashids(hid, content.PostId)
	if err != nil {
		return 0, err
	}
	return content.PostId, nil
}

// UpdatePostHashids 更新文章 Hashids
func (service Service) UpdatePostHashids(postSlug string, postId uint64) error {
	if err := service.db.Table(model.TPostsTable).
		Where("post_id", postId).
		Update("post_slug", postSlug).Error; err != nil {
		return errors.New("更新失败: " + err.Error())
	}
	return nil
}

// UpdatePost 更新文章
func (service Service) UpdatePost(obj model.Post) error {
	if err := service.db.Table(model.TPostsTable).Save(&obj).Error; err != nil {
		slog.Errorf("更新文章失败: %s", err.Error())
		return errors.New("更新失败: " + err.Error())
	}
	return nil
}

// PublishArticle 发布/取消发布文章
func (service Service) PublishArticle(obj model.Post) error {
	if err := service.db.Table(model.TPostsTable).
		Where("post_id", obj.PostId).
		Update("is_published", obj.IsPublished).Error; err != nil {
		slog.Errorf("更新发布状态失败: %s", err.Error())
		return errors.New("更新失败: " + err.Error())
	}
	return nil
}

// TopPost 置顶/取消置顶文章
func (service Service) TopPost(postId int64, top uint8) error {
	if err := service.db.Table(model.TPostsTable).
		Where("post_id", postId).
		Update("top", top).Error; err != nil {
		slog.Errorf("更新置顶状态失败: %s", err.Error())
		return errors.New("更新失败: " + err.Error())
	}
	return nil
}

// DeletePost 删除文章（软删除）
func (service Service) DeletePost(postId int) error {
	if err := service.db.Table(model.TPostsTable).
		Where("post_id", postId).
		Update("is_deleted", 1).Error; err != nil {
		return errors.New("删除失败: " + err.Error())
	}
	return nil
}

// UpdatePostCoverImag 更新文章封面
func (service Service) UpdatePostCoverImag(postId int) error {
	if err := service.db.Table(model.TPostsTable).
		Where("post_id", postId).
		Update("cover_image", pkg.GetPixabayImage()).Error; err != nil {
		slog.Errorf("更新封面失败: %s", err.Error())
		return errors.New("更新封面失败: " + err.Error())
	}
	return nil
}

// UpdatePostAllCoverImag 批量更新所有文章封面
func (service Service) UpdatePostAllCoverImag() error {
	var content []model.Post

	err := NewQueryBuilder(service.db).
		WithPublished().
		WithNotDeleted().
		Find(&content)

	if err != nil {
		return errors.New("查询文章失败: " + err.Error())
	}

	for _, post := range content {
		service.db.Table(model.TPostsTable).
			Where("post_id", post.PostId).
			Update("cover_image", pkg.GetPixabayImage())
	}

	return nil
}

// GetList 获取文章列表（后台管理）
func (service Service) GetList(req request.PostRequest) (*page.Info, error) {
	var content []model.Post
	var total int64
	pageNum := req.PageNum
	pageSize := req.PageSize

	qb := NewQueryBuilder(service.db).
		Select("post.*, pc.category_name").
		WithNotDeleted().
		WithTitle(req.Title).
		WithPublishedStatus(req.Published).
		WithCategoryId(req.CategoryId)

	// 关联分类表（使用 join 而不是 left join，因为后台需要显示分类）
	qb.DB = qb.DB.Joins("JOIN category pc ON post.category_id = pc.category_id")

	if err := qb.Count(&total); err != nil {
		return nil, errors.New("统计文章失败: " + err.Error())
	}
	if err := qb.
		Order("post_id DESC").
		WithPagination(pageNum, pageSize).
		Find(&content); err != nil {
		return nil, errors.New("查询文章失败: " + err.Error())
	}

	return page.PaginationInfo(content, pageNum, pageSize, int(total)), nil
}
