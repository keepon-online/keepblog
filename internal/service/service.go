package service

import (
	"gitee.com/jieepre/keepblog/internal/service/about"
	"gitee.com/jieepre/keepblog/internal/service/category"
	"gitee.com/jieepre/keepblog/internal/service/dashboard"
	"gitee.com/jieepre/keepblog/internal/service/link"
	"gitee.com/jieepre/keepblog/internal/service/music"
	"gitee.com/jieepre/keepblog/internal/service/notice"
	"gitee.com/jieepre/keepblog/internal/service/post"
	"gitee.com/jieepre/keepblog/internal/service/sidebar"
	"gitee.com/jieepre/keepblog/internal/service/system"
	"gitee.com/jieepre/keepblog/internal/service/tag"
	"gitee.com/jieepre/keepblog/internal/service/website"
	"gorm.io/gorm"
)

type AppService struct {
	PostService     *post.Service
	CategoryService *category.Service
	TagService      *tag.Service
	SidebarService  *sidebar.Service
	AboutService    *about.Service
	SystemService   *system.Service
	LinkService     *link.Service
	WebSiteService  *website.Service
	Dashboard       *dashboard.Service
	MusicService    *music.Service
	NoticeService   *notice.Service
}

// InitAppService 装配应用服务。数据库连接由调用方传入，
// 服务层不再依赖全局变量（global.GORM 仅存于 app.Initialize）。
func InitAppService(db *gorm.DB) *AppService {
	return &AppService{
		PostService:     post.NewPostService(db),
		CategoryService: category.NewCategoryService(db),
		TagService:      tag.NewTagService(db),
		SidebarService:  sidebar.NewSidebarService(db),
		AboutService:    about.NewAboutService(db),
		SystemService:   system.NewSystemService(db),
		LinkService:     link.NewLinkService(db),
		WebSiteService:  website.NewWebSiteService(db),
		Dashboard:       dashboard.NewDashboardService(db),
		MusicService:    music.NewMusicService(db),
		NoticeService:   notice.NewNoticeService(db),
	}
}
