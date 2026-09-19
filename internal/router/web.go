package router

import (
	"io/fs"
	"net/http"
	"time"

	"gitee.com/jieepre/keepblog/api/web/about"
	"gitee.com/jieepre/keepblog/api/web/archive"
	"gitee.com/jieepre/keepblog/api/web/category"
	"gitee.com/jieepre/keepblog/api/web/home"
	"gitee.com/jieepre/keepblog/api/web/link"
	"gitee.com/jieepre/keepblog/api/web/music"
	"gitee.com/jieepre/keepblog/api/web/post"
	"gitee.com/jieepre/keepblog/api/web/tags"
	"gitee.com/jieepre/keepblog/internal/middleware"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/static"
	"github.com/gin-gonic/gin"
)

// RegisterWebRouter 注册前台路由。
// group 是挂有限流/统计/页面缓存中间件的前台分组，页面路由必须注册在
// 它的子分组上才会继承中间件；静态资源、Feed 与音乐 API 因各自的缓存
// 策略不同，仍直接注册在 Engine 上。
func RegisterWebRouter(ctx *core.Context, group *gin.RouterGroup) {
	staticRouter(ctx)
	feedRouter(ctx)
	aboutRouter(ctx, group)
	tagWebRouter(ctx, group)
	linkRouter(ctx, group)
	archivesRouter(ctx, group)
	categoriesRouter(ctx, group)
	postsRouter(ctx, group)
	homeRouter(ctx, group)
	musicWebRouter(ctx)
}

func staticRouter(ctx *core.Context) {
	// 静态资源统一走缓存中间件：embed 内容随发版才变，配 1 年 immutable 强缓存。
	// 注意：StaticFS/StaticFileFS 必须注册在带中间件的 group 上，中间件才生效。
	staticGroup := ctx.Engine.Group("/", middleware.StaticCacheMiddleware())

	cssEmbed, _ := fs.Sub(static.Static, "css")
	jsEmbed, _ := fs.Sub(static.Static, "js")
	imagesEmbed, _ := fs.Sub(static.Static, "images")
	pluginsEmbed, _ := fs.Sub(static.Static, "plugins")
	staticGroup.StaticFS("/css", http.FS(cssEmbed))
	staticGroup.StaticFS("/js", http.FS(jsEmbed))
	staticGroup.StaticFS("/images", http.FS(imagesEmbed))
	staticGroup.StaticFS("/plugins", http.FS(pluginsEmbed))
	// robots.txt 改为 feedRouter 里的动态 handler（需注入站点 sitemap 地址）
	staticGroup.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(static.Favicon))
}

// feedRouter 注册 RSS / sitemap / robots 订阅与索引路由。
// 这些内容来自 DB（站点信息 + 已发布文章），需动态生成，因此不挂在 staticGroup。
// Feed 注册在 Engine 上（不经过 webGroup），因此显式为每条路由配置缓存。
func feedRouter(ctx *core.Context) {
	handler := post.Handler{Context: ctx}
	group := ctx.Engine.Group("/")
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/rss.xml",
			Handler:    handler.Feed,
			Middleware: []gin.HandlerFunc{middleware.CacheMiddleware(30 * time.Minute)},
		},
		{
			Method:     http.MethodGet,
			Path:       "/sitemap.xml",
			Handler:    handler.Sitemap,
			Middleware: []gin.HandlerFunc{middleware.CacheMiddleware(1 * time.Hour)},
		},
		{
			Method:     http.MethodGet,
			Path:       "/robots.txt",
			Handler:    handler.Robots,
			Middleware: []gin.HandlerFunc{middleware.CacheMiddleware(1 * time.Hour)},
		},
	}
	RegisterRouter(group, routes)
}

func homeRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := home.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "首页",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/",
			Handler:    handler.Home,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/page/:page",
			Handler:    handler.Home,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/search/:keyword",
			Handler:    handler.Search,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodGet,
			Path:       "/daily",
			Handler:    handler.Daily,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func aboutRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := about.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "关于",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/about",
			Handler:    handler.About,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func tagWebRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := tags.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "标签",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/tags",
			Handler:    handler.Tags,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/tags/:tag",
			Handler:    handler.TagsPage,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/tags/:tag/page/:page",
			Handler:    handler.TagsPage,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func linkRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := link.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "友情链接",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/link",
			Handler:    handler.Links,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func archivesRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := archive.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "归档",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/archives",
			Handler:    handler.Archives,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/archives/page/:page",
			Handler:    handler.Archives,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/archives/:year",
			Handler:    handler.ArchivesInfo,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/archives/:year/:month",
			Handler:    handler.ArchivesInfo,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func categoriesRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := category.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "分类",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/categories",
			Handler:    handler.Categories,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/categories/:category",
			Handler:    handler.CategoriesPage,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/categories/:category/page/:page",
			Handler:    handler.CategoriesPage,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func postsRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := post.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "查看文章",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/post/:hashids",
			Handler:    handler.Post,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

// musicWebRouter 音乐列表是公开 JSON API,不走页面缓存,注册在 Engine 上。
func musicWebRouter(ctx *core.Context) {
	handler := music.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "音乐",
		Prefix: "/api",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/music/list",
			Handler:    handler.GetMusicList,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}
