package service

import (
	"gitee.com/jieepre/go-site/internal/service/about"
	"gitee.com/jieepre/go-site/internal/service/category"
	"gitee.com/jieepre/go-site/internal/service/dashboard"
	"gitee.com/jieepre/go-site/internal/service/link"
	"gitee.com/jieepre/go-site/internal/service/music"
	"gitee.com/jieepre/go-site/internal/service/post"
	"gitee.com/jieepre/go-site/internal/service/sidebar"
	"gitee.com/jieepre/go-site/internal/service/system"
	"gitee.com/jieepre/go-site/internal/service/tag"
	"gitee.com/jieepre/go-site/internal/service/website"
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
}

func InitAppService() *AppService {
	return &AppService{
		PostService:     post.NewPostService(),
		CategoryService: category.NewCategoryService(),
		TagService:      tag.NewTagService(),
		SidebarService:  sidebar.NewSidebarService(),
		AboutService:    about.NewAboutService(),
		SystemService:   system.NewSystemService(),
		LinkService:     link.NewLinkService(),
		WebSiteService:  website.NewWebSiteService(),
		Dashboard:       dashboard.NewDashboardService(),
		MusicService:    music.NewMusicService(),
	}
}
