package sidebar

import (
	"time"

	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/system"
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
		WebInfo:         service.WebInfo(),
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

// WebInfo 获取网站资讯详细数据
func (service Service) WebInfo() model.WebInfo {
	var webInfo model.WebInfo

	// 文章数目
	global.GORM.Table(model.TPostsTable).
		Where("is_published", 1).
		Where("is_deleted", 0).
		Count(&webInfo.PostCount)

	// 总字数 - 使用 Raw SQL 保证聚合函数正确执行
	var totalWordCount struct {
		Total int64 `gorm:"column:total"`
	}
	global.GORM.Raw("SELECT COALESCE(SUM(word_count), 0) as total FROM post WHERE is_published = 1 AND is_deleted = 0").Scan(&totalWordCount)
	webInfo.TotalWordCount = totalWordCount.Total

	// 从网站配置读取创建日期
	var siteConfig system.WebSite
	global.GORM.Model(&system.WebSite{}).First(&siteConfig)

	// 解析网站创建日期，默认为 2023-05-20
	siteStartDateStr := siteConfig.SiteStartDate
	if siteStartDateStr == "" {
		siteStartDateStr = "2023-05-20"
	}
	siteStartDate, err := time.ParseInLocation("2006-01-02", siteStartDateStr, time.Local)
	if err != nil {
		siteStartDate = time.Date(2023, 5, 20, 0, 0, 0, 0, time.Local)
	}

	webInfo.SiteStartDate = siteStartDate.Format("2006-01-02")
	webInfo.RuntimeDays = int64(time.Since(siteStartDate).Hours() / 24)

	// 最后更新时间 - 使用 Raw SQL 获取最新的更新时间
	var lastUpdate struct {
		LastTime uint64 `gorm:"column:last_time"`
	}
	global.GORM.Raw("SELECT COALESCE(MAX(last_modified_time), 0) as last_time FROM post WHERE is_published = 1 AND is_deleted = 0").Scan(&lastUpdate)

	if lastUpdate.LastTime > 0 {
		webInfo.LastUpdateTime = time.Unix(int64(lastUpdate.LastTime), 0).Format("2006年1月2日")
	} else {
		// 如果没有最后更新时间，尝试使用最新的发布时间
		var pubTime struct {
			PubTime uint64 `gorm:"column:pub_time"`
		}
		global.GORM.Raw("SELECT COALESCE(MAX(pub_time), 0) as pub_time FROM post WHERE is_published = 1 AND is_deleted = 0").Scan(&pubTime)
		if pubTime.PubTime > 0 {
			webInfo.LastUpdateTime = time.Unix(int64(pubTime.PubTime), 0).Format("2006年1月2日")
		} else {
			webInfo.LastUpdateTime = "暂无文章"
		}
	}

	return webInfo
}
