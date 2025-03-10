package sidebar

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/service/category"
	"gitee.com/jieepre/go-site/internal/service/post"
	"gitee.com/jieepre/go-site/internal/service/tag"
)

type Service struct {
}

func NewSidebarService() *Service {
	return &Service{}
}

func (service Service) Sidebar() *model.Sidebar {
	sidebar := model.Sidebar{
		Category:        service.Categories(),
		Tag:             service.Tag(),
		CardInfo:        service.CardInfo(),
		LatestPosts:     service.LatestPosts(),
		SidebarArchives: service.SidebarArchives(),
	}
	return &sidebar
}

func (service Service) CardInfo() model.CardInfo {
	var cardCount model.CardInfo
	global.GORM.Table(model.TPostsTable).
		Select("COUNT(DISTINCT t.tag_id )  AS tag").
		Where("post.is_published", 1).
		Where("post.is_deleted", 0).
		Joins("LEFT JOIN post_tag pt on post.post_id = pt.post_id").
		Joins("LEFT JOIN tag t on pt.tag_id = t.tag_id").
		Count(&cardCount.Tag)
	global.GORM.Table(model.TPostsTable).
		Select("COUNT(DISTINCT category_id)  AS category").
		Where("is_published", 1).
		Where("is_deleted", 0).
		Count(&cardCount.Category)
	global.GORM.Table(model.TPostsTable).
		Where("is_published", 1).
		Where("is_deleted", 0).Count(&cardCount.Post)
	return cardCount
}

func (service Service) Tag() []model.Tag {
	tagService := tag.NewTagService()
	tags, _ := tagService.GetTags()
	return tags
}

func (service Service) LatestPosts() []model.LatestPosts {
	latestPosts, _ := post.NewPostService().GetLatestPosts()
	return latestPosts
}
func (service Service) Categories() []model.CategoryCount {

	categoryService := category.NewCategoryService()
	categories, _ := categoryService.GetCategories()
	return categories
}

func (service Service) SidebarArchives() []model.SidebarArchives {
	var years []string
	global.GORM.Table(model.TPostsTable).Raw(`SELECT strftime( '%Y/%m', pub_time, 'unixepoch' ) year FROM post WHERE is_published = 1 AND is_deleted = 0 GROUP BY year`).Scan(&years)
	m := make(map[string]int64)
	sidebarArchives := make([]model.SidebarArchives, 0)
	for _, year := range years {
		var total int64
		global.GORM.Table(model.TPostsTable).Raw(" SELECT count(0) from post WHERE  is_published = 1 AND is_deleted = 0 AND strftime( '%Y/%m', pub_time, 'unixepoch' ) =?", year).Scan(&total)
		m[year] = total
		sidebarArchive := model.SidebarArchives{
			Year:  year,
			Total: total,
		}
		sidebarArchives = append(sidebarArchives, sidebarArchive)
	}
	return sidebarArchives
}
