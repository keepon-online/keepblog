package post

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"github.com/pkg/errors"
)

type Service struct {
}

func NewPostService() *Service {
	return &Service{}
}

func (service Service) GetPost(id int) (*model.Post, error) {
	var postInfo model.Post

	if err := global.GORM.Table(model.TPostsTable).
		Select("post.*,category.category_name,t.tag_name").
		Where("post.post_id", id).
		Where("post.is_published", 1).
		Where("post.is_deleted", 0).
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
		Where("post.is_published", 1).
		Where("post.is_deleted", 0).
		Joins("LEFT JOIN post_tag pt on post.post_id = pt.post_id").
		Joins("LEFT JOIN tag t on pt.tag_id = t.tag_id").
		Scan(&m).Error; err != nil {
		return nil, errors.New("找不到记录")
	}

	postInfo.Tags = m
	_ = service.UpdatePostReadCount(postInfo.PostId, postInfo.ReadCount)

	return &postInfo, nil
}

// GetLatestPosts 最新文章
func (service Service) GetLatestPosts() ([]model.LatestPosts, error) {
	var latestPosts []model.LatestPosts
	err := global.GORM.Table(model.TPostsTable).
		Select("post.title,post.author,post.pub_time,post.cover_image, post.post_slug,post.pub_time,category.category_name,category.category_id").
		Joins("LEFT JOIN category  ON  post.category_id=category.category_id").
		Where("is_published", 1).
		Where("is_deleted", 0).
		Limit(5).Find(&latestPosts).Order("post.create_time desc").Error

	if err != nil {
		return nil, errors.New("查询失败" + err.Error())
	}

	return latestPosts, nil
}

func (service Service) GetCoverPosts(pageNum int) ([]model.LatestPosts, int64, error) {
	var coverPosts []model.LatestPosts
	var total int64
	db := global.GORM.Table(model.TPostsTable).
		Select("post.title,post.author,post.pub_time,post.cover_image,post.top, post.post_slug,post.pub_time,category.category_name,category.category_id,t.tag_name").
		Joins("LEFT JOIN category ON  post.category_id=category.category_id").
		Joins("LEFT JOIN post_tag pt on post.post_id = pt.post_id").
		Joins("LEFT JOIN tag t on pt.tag_id = t.tag_id").
		Where("is_published", 1).
		Where("is_deleted", 0)
	err := db.Count(&total).Error
	err = db.Limit(10).Offset((pageNum - 1) * 10).
		Order(" post.top desc,post.create_time desc ").
		Find(&coverPosts).Error

	if err != nil {
		return nil, 0, errors.New("查询失败" + err.Error())
	}

	return coverPosts, total, nil
}

func (service Service) Total() (total int64) {
	global.GORM.Table(model.TPostsTable).Where("is_published", 1).Where("is_deleted", 0).Count(&total)
	return
}

func (service Service) UpdatePostReadCount(postId uint64, readCount uint32) error {
	err := global.GORM.Table(model.TPostsTable).Where("post_id", postId).Update("read_count", readCount+1).Error
	if err != nil {
		return errors.New("更新失败" + err.Error())
	}
	return nil
}

func (service Service) GetArchivePosts(year, month string) (*model.ArchivesPosts, error) {
	var years []string
	global.GORM.Table(model.TPostsTable).Raw(`SELECT strftime( '%Y-%m', pub_time, 'unixepoch' ) year FROM post where is_published = 1 AND is_deleted = 0 GROUP BY year`).Scan(&years)
	m := make(map[string][]model.ArchivePosts, 0)

	archives := model.ArchivesPosts{}
	if year != "" && month != "" {
		posts := make([]model.ArchivePosts, 0)
		global.GORM.Table(model.TPostsTable).Raw("SELECT title,post_slug,pub_time,cover_image FROM post WHERE is_published = 1 AND is_deleted = 0 AND  strftime('%Y/%m', pub_time, 'unixepoch' ) =?", year+"/"+month).Scan(&posts)
		m[year+"-"+month] = posts
		archives.Archives = m
		return &archives, nil
	}

	for _, year := range years {
		posts := make([]model.ArchivePosts, 0)
		global.GORM.Table(model.TPostsTable).Raw("SELECT title,post_slug,pub_time,cover_image FROM post WHERE is_published = 1 AND is_deleted = 0 AND strftime( '%Y-%m', pub_time, 'unixepoch' ) =?", year).Scan(&posts)
		m[year] = posts
	}

	archives.Archives = m
	return &archives, nil
}

func (service Service) GetPostsByCategory(category string, pageNum int) ([]model.TagCategoryPosts, int64, error) {
	pageSize := 10
	var total int64
	posts := make([]model.TagCategoryPosts, 0)

	tx := global.GORM.Table(model.TPostsTable).
		Select("post.title, post.post_slug, post.pub_time,post.cover_image").
		Joins("LEFT JOIN category  ON  post.category_id=category.category_id").
		Where("category.category_name", category).
		Where("post.is_published", 1).
		Where("post.is_deleted", 0)
	tx.Count(&total)
	tx.Limit(pageSize).Offset((pageNum - 1) * pageSize).Scan(&posts)

	return posts, total, nil
}

func (service Service) GetPostsByTag(tagName string, pageNum int) ([]model.TagCategoryPosts, int64, error) {
	pageSize := 10
	var total int64
	posts := make([]model.TagCategoryPosts, 0)

	tx := global.GORM.Table(model.TPostsTable).
		Select("post.title, post.post_slug, post.pub_time,post.cover_image").
		Joins("LEFT JOIN post_tag pt on post.post_id = pt.post_id").
		Joins("LEFT JOIN tag t on pt.tag_id = t.tag_id").
		Where("post.is_published", 1).
		Where("post.is_deleted", 0).
		Where("t.tag_name", tagName)
	tx.Count(&total)
	tx.Limit(pageSize).Offset((pageNum - 1) * pageSize).Scan(&posts)
	return posts, total, nil
}

func (service Service) GetPostsByCategoryId(id int64) ([]model.Post, error) {
	posts := make([]model.Post, 0)
	tx := global.GORM.Table(model.TPostsTable).Where("category_id", id).Find(&posts)

	if tx.Error != nil {
		return nil, errors.New("查询失败")
	}

	return posts, nil
}

func (service Service) GetPostsArchive(num int64) ([]map[string]any, int64, error) {
	var archives []struct {
		Year  string
		Month string
		Count int
	}
	global.GORM.Raw("SELECT strftime('%Y', create_time, 'unixepoch' ) AS year, strftime('%m', create_time,'unixepoch') AS month, COUNT(*) AS count FROM post GROUP BY year, month ORDER BY year DESC, month DESC").Scan(&archives)

	// 生成归档数据
	archiveData := make([]map[string]any, 0)
	for i, archive := range archives {
		// 查询归档时间段内的文章
		var articles []model.ArchivePosts
		global.GORM.Table("post").Where("strftime('%Y', create_time, 'unixepoch') = ? AND strftime('%m', create_time, 'unixepoch') = ?", archive.Year, archive.Month).Order("create_time DESC").Find(&articles)
		// 生成文章列表数据
		articleData := make([]map[string]interface{}, 0)
		for j, article := range articles {
			articleData[j] = map[string]interface{}{
				"title":      article.Title,
				"created_at": article.PubTime,
			}
		}

		// 生成归档数据
		archiveData[i] = map[string]any{
			"year":     archive.Year,
			"month":    archive.Month,
			"count":    archive.Count,
			"articles": articleData,
		}
	}

	return archiveData, 1, nil
}

func (service Service) Search(keyword string) (posts []model.SearchPost, err error) {
	err = global.GORM.Table(model.TPostsTable).Where("title like ? and is_published=1 and is_deleted=0 ", "%"+keyword+"%").Find(&posts).Error
	return
}
