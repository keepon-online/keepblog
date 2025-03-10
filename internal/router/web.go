package router

import (
	"gitee.com/jieepre/go-site/api/web/about"
	"gitee.com/jieepre/go-site/api/web/archive"
	"gitee.com/jieepre/go-site/api/web/category"
	"gitee.com/jieepre/go-site/api/web/home"
	"gitee.com/jieepre/go-site/api/web/link"
	"gitee.com/jieepre/go-site/api/web/post"
	"gitee.com/jieepre/go-site/api/web/tags"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

func RegisterWebRouter(ctx *core.Context) {
	staticRouter(ctx)
	aboutRouter(ctx)
	tagWebRouter(ctx)
	linkRouter(ctx)
	archivesRouter(ctx)
	categoriesRouter(ctx)
	postsRouter(ctx)
	homeRouter(ctx)
}

func staticRouter(ctx *core.Context) {
	web := ctx.Engine
	cssEmbed, _ := fs.Sub(static.Static, "css")
	jsEmbed, _ := fs.Sub(static.Static, "js")
	imagesEmbed, _ := fs.Sub(static.Static, "images")
	pluginsEmbed, _ := fs.Sub(static.Static, "plugins")
	web.StaticFS("/css", http.FS(cssEmbed))
	web.StaticFS("/js", http.FS(jsEmbed))
	web.StaticFS("/images", http.FS(imagesEmbed))
	web.StaticFS("/plugins", http.FS(pluginsEmbed))
	web.StaticFileFS("/robots.txt", "robots.txt", http.FS(static.Robots))
	web.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(static.Favicon))
}

func homeRouter(ctx *core.Context) {
	handler := home.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "首页",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
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

func aboutRouter(ctx *core.Context) {
	handler := about.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "关于",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
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

func tagWebRouter(ctx *core.Context) {
	handler := tags.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "标签",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
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

func linkRouter(ctx *core.Context) {
	handler := link.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "友情链接",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
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

func archivesRouter(ctx *core.Context) {
	handler := archive.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "归档",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
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

func categoriesRouter(ctx *core.Context) {
	handler := category.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "分类",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
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

func postsRouter(ctx *core.Context) {
	handler := post.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "查看文章",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
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
