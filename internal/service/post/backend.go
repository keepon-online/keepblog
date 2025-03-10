package post

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/request"
	"gitee.com/jieepre/go-site/pkg"
	"gitee.com/jieepre/go-site/pkg/hash"
	"gitee.com/jieepre/go-site/pkg/page"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (service Service) SavePost(content model.Post) (error, uint64) {
	if err := global.GORM.Table(model.TPostsTable).Create(&content).Error; err != nil {
		return errors.New("保存失败 " + err.Error()), content.PostId
	}
	hid, _ := hash.New().HashidsEncode([]int{int(content.PostId)})
	err := service.UpdatePostHashids(hid, content.PostId)
	if err != nil {
		return err, 0
	}
	return nil, content.PostId
}
func (service Service) GetPostDetail(id int) (*model.Post, error) {
	var postInfo model.Post

	if err := global.GORM.Table(model.TPostsTable).
		Select("post.*,category.category_name,t.tag_name").
		Where("post.post_id", id).
		Joins("LEFT JOIN category ON post.category_id=category.category_id").
		Joins("LEFT JOIN post_tag pt on post.post_id = pt.post_id").
		Joins("LEFT JOIN tag t on pt.tag_id = t.tag_id").
		Scan(&postInfo).Error; err != nil {
		return nil, errors.New("找不到记录")
	}

	m := make([]string, 0)
	if err := global.GORM.Table(model.TPostsTable).
		Select("t.tag_name").
		Where("post.post_id", id).
		Joins("LEFT JOIN post_tag pt on post.post_id = pt.post_id").
		Joins("LEFT JOIN tag t on pt.tag_id = t.tag_id").
		Scan(&m).Error; err != nil {
		return nil, errors.New("找不到记录")
	}
	postInfo.Tags = m
	_ = service.UpdatePostReadCount(postInfo.PostId, postInfo.ReadCount)

	return &postInfo, nil
}

func (service Service) UpdatePostHashids(postSlug string, postId uint64) error {
	if err := global.GORM.Table(model.TPostsTable).
		Where("post_id", postId).
		Update("post_slug", postSlug).Error; err != nil {
		return errors.New("update post error " + err.Error())
	}
	return nil
}
func (service Service) UpdatePost(obj model.Post) error {
	if err := global.GORM.Table(model.TPostsTable).Save(&obj).Error; err != nil {
		slog.Errorf("update post error %s", err.Error())
		return errors.New("update post error " + err.Error())
	}
	return nil
}

func (service Service) PublishArticle(obj model.Post) error {
	if err := global.GORM.Table(model.TPostsTable).Where("post_id", obj.PostId).Update("is_published", obj.IsPublished).Error; err != nil {
		slog.Errorf("update post error:%s", err.Error())
		return errors.New("update post error " + err.Error())
	}
	return nil
}

func (service Service) TopPost(postId int64, top uint8) error {
	if err := global.GORM.Table(model.TPostsTable).Where("post_id", postId).Update("top", top).Error; err != nil {
		slog.Errorf("update post error %s", err.Error())
		return errors.New("update post error " + err.Error())
	}
	return nil
}

func (service Service) DeletePost(postId int) error {
	if err := global.GORM.Table(model.TPostsTable).Where("post_id", postId).Update("is_deleted", 1).Error; err != nil {
		return errors.New("delete post error " + err.Error())
	}
	return nil
}

func (service Service) UpdatePostCoverImag(postId int) error {
	if err := global.GORM.Table(model.TPostsTable).Where("post_id", postId).Update("cover_image", pkg.GetPixabayImage()).Error; err != nil {
		slog.Errorf("update post error:%s", err.Error())
		return errors.New("update cover error " + err.Error())
	}
	return nil
}

func (service Service) UpdatePostAllCoverImag() error {
	var content []model.Post
	global.GORM.Table(model.TPostsTable).Where("is_published=1 and is_deleted=0").Find(&content)
	for _, post := range content {
		global.GORM.Table(model.TPostsTable).Where("post_id", post.PostId).Update("cover_image", pkg.GetPixabayImage())
	}

	return nil
}

func (service Service) GetList(req request.PostRequest) *page.Info {
	var content []model.Post
	var total int64
	pageNum := req.PageNum
	pageSize := req.PageSize
	global.GORM.Table(model.TPostsTable).
		Select("post.*,pc.category_name").
		Scopes(IsDeleted, postTitle(req.Title), postPublished(req.Published), postCategory(req.CategoryId)).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Order("post_id desc").
		Joins("join category pc on post.category_id = pc.category_id ").
		Find(&content)

	global.GORM.Table(model.TPostsTable).
		Scopes(IsDeleted, postTitle(req.Title), postPublished(req.Published), postCategory(req.CategoryId)).
		Joins("join category pc on post.category_id = pc.category_id ").
		Count(&total)
	bInfo := page.PaginationInfo(content, pageNum, pageSize, int(total))
	return bInfo
}
func IsDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("post.is_deleted", 0)
}
func postTitle(title string) func(db *gorm.DB) *gorm.DB {
	if title == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("post.title like ?", "%"+title+"%")
	}
}
func postCategory(categoryId *uint32) func(db *gorm.DB) *gorm.DB {
	if categoryId == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("post.category_id", categoryId)
	}
}
func postPublished(published *uint8) func(db *gorm.DB) *gorm.DB {
	if published == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("post.is_published", published)
	}
}
