package post

import (
	"errors"
	"time"

	"github.com/gookit/slog"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/model/request"
	"gitee.com/jieepre/keepblog/pkg"
	"gitee.com/jieepre/keepblog/pkg/hash"
	"gitee.com/jieepre/keepblog/pkg/md"
	"gitee.com/jieepre/keepblog/pkg/page"
)

// SavePost 保存文章
func (service Service) SavePost(content model.Post) (uint64, error) {
	// pub_time 不再依赖 gorm autoUpdateTime（Save 会覆盖非零值）：
	// 新建未指定发布时间时默认当前时间，指定（定时发布）则原样保留。
	if content.PubTime == 0 {
		content.PubTime = uint64(time.Now().Unix())
	}
	content.WordCount = uint32(md.CountWords([]byte(content.PostContent)))
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
	var old model.Post
	if err := service.db.Table(model.TPostsTable).
		Where("post_id", obj.PostId).
		Take(&old).Error; err == nil {
		// 1. pub_time 显式管理：请求未指定时保留原发布时间，避免 Save 把发布时间
		// 重置为编辑时刻（定时发布的未来时间也在此保护之下）。
		if obj.PubTime == 0 {
			obj.PubTime = old.PubTime
		}
		// 2. 保护 post_slug：若请求未传或为空，保留旧 slug；若旧 slug 亦为空，则基于 PostId 重新生成
		if obj.PostSlug == "" {
			if old.PostSlug != "" {
				obj.PostSlug = old.PostSlug
			} else {
				hid, _ := hash.New().HashidsEncode([]int{int(obj.PostId)})
				obj.PostSlug = hid
			}
		}
		// 3. 保护发布状态：普通编辑请求通常不携带 published 变更，避免被结构体零值(0)误改回未发布状态
		if obj.IsPublished == 0 && old.IsPublished != 0 {
			obj.IsPublished = old.IsPublished
		}
		// 4. 保护创建时间与阅读量，避免被 0 覆盖
		if obj.CreateTime == 0 {
			obj.CreateTime = old.CreateTime
		}
		if obj.ReadCount == 0 {
			obj.ReadCount = old.ReadCount
		}
		// 5. 若 Author 未传，保留旧 Author
		if obj.Author == "" {
			obj.Author = old.Author
		}
	} else if obj.PostSlug == "" && obj.PostId != 0 {
		hid, _ := hash.New().HashidsEncode([]int{int(obj.PostId)})
		obj.PostSlug = hid
	}

	obj.WordCount = uint32(md.CountWords([]byte(obj.PostContent)))
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
