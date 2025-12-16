package router

import (
	"net/http"

	"gitee.com/jieepre/go-site/api/admin/about"
	"gitee.com/jieepre/go-site/api/admin/category"
	"gitee.com/jieepre/go-site/api/admin/common"
	"gitee.com/jieepre/go-site/api/admin/dashboard"
	"gitee.com/jieepre/go-site/api/admin/link"
	"gitee.com/jieepre/go-site/api/admin/login"
	"gitee.com/jieepre/go-site/api/admin/monitor"
	"gitee.com/jieepre/go-site/api/admin/music"
	"gitee.com/jieepre/go-site/api/admin/post"
	"gitee.com/jieepre/go-site/api/admin/system"
	"gitee.com/jieepre/go-site/api/admin/tags"
	"gitee.com/jieepre/go-site/api/admin/website"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(ctx *core.Context) {
	userRouter(ctx)
	monitorRouter(ctx)
	categoryRouter(ctx)
	tagAdminRouter(ctx)
	postRouter(ctx)
	commonRouter(ctx)
	aboutBackendRouter(ctx)
	linkBackendRouter(ctx)
	websiteBackendRouter(ctx)
	logsRouter(ctx)
	dashboardRouter(ctx)
	musicRouter(ctx)
}

func userRouter(ctx *core.Context) {
	handler := login.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "系统管理",
		Prefix: "/",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/api/login",
			Handler:    handler.Login,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodPost,
			Path:       "/api/logout",
			Handler:    handler.Logout,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodPost,
			Path:       "/api/change-password",
			Handler:    handler.ChangePassword,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodPost,
			Path:       "/api/refreshToken",
			Handler:    handler.RefreshToken,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/api/getAsyncRoutes",
			Handler:    handler.GetAsyncRoutes,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/api/getInfo",
			Handler:    handler.GetInfo,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func monitorRouter(ctx *core.Context) {
	handler := monitor.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "监控",
		Prefix: "/api/monitor",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/server",
			Handler:    handler.Monitor,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodGet,
			Path:       "/realtime",
			Handler:    handler.GetRealtime,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodGet,
			Path:       "/general",
			Handler:    handler.General,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/loadavg",
			Handler:    handler.Loadavg,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/ram",
			Handler:    handler.RAM,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/cpu",
			Handler:    handler.CPU,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/net",
			Handler:    handler.Net,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/diskUsage",
			Handler:    handler.DiskUsage,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/diskIOStat",
			Handler:    handler.DiskIOStat,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodGet,
			Path:       "/base/os",
			Handler:    handler.LoadDashboardOsInfo,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodGet,
			Path:       "/base/:ioOption/:netOption",
			Handler:    handler.LoadDashboardBaseInfo,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodGet,
			Path:       "/current/:ioOption/:netOption",
			Handler:    handler.LoadDashboardCurrentInfo,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func categoryRouter(ctx *core.Context) {
	handler := category.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "分类管理",
		Prefix: "/api/v1/site/category",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveCategory,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateCategory,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateCategoryState,
			Middleware: nil,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:categoryId",
			Handler:    handler.DeleteCategory,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:categoryId",
			Handler:    handler.DetailCategory,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetCategoryList,
			Middleware: nil,
		},
	}
	RegisterRouter(group, routes)
}

func postRouter(ctx *core.Context) {
	handler := post.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "测试分组",
		Prefix: "/api/v1/site/post",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SavePost,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdatePost,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-publish",
			Handler:    handler.PublishPost,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-top",
			Handler:    handler.TopPost,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-cover/:postId",
			Handler:    handler.UpdatePostCoverImag,
			Middleware: nil,
		}, {
			Method:     http.MethodPut,
			Path:       "/update-cover",
			Handler:    handler.UpdatePostAllCoverImag,
			Middleware: nil,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:postId",
			Handler:    handler.DeletePost,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:postId",
			Handler:    handler.DetailPost,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetList,
			Middleware: nil,
		},
	}
	RegisterRouter(group, routes)
}

func tagAdminRouter(ctx *core.Context) {
	handler := tags.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "标签管理",
		Prefix: "/api/v1/site/tags",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetTags,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func commonRouter(ctx *core.Context) {
	handler := common.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "公用",
		Prefix: "/api/v1",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/upload/images",
			Handler:    handler.UploadImage,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func aboutBackendRouter(ctx *core.Context) {
	handler := about.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "关于",
		Prefix: "/api/v1/site/about",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.Save,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/detail",
			Handler:    handler.Detail,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func linkBackendRouter(ctx *core.Context) {
	handler := link.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "友链管理",
		Prefix: "/api/v1/site/link",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveLink,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateLink,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateLinkState,
			Middleware: nil,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteLink,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:id",
			Handler:    handler.DetailLink,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetLinkList,
			Middleware: nil,
		},
	}
	RegisterRouter(group, routes)
}

func websiteBackendRouter(ctx *core.Context) {
	handler := website.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "网站",
		Prefix: "/api/v1/site/web",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.SaveWebSite,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/detail",
			Handler:    handler.GetWebSite,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func logsRouter(ctx *core.Context) {
	handler := system.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "网站",
		Prefix: "/api/v1/site/logs",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/logon",
			Handler:    handler.LoginLog,
			Middleware: []gin.HandlerFunc{},
		}, {
			Method:     http.MethodGet,
			Path:       "/access",
			Handler:    handler.AccessLog,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}
func dashboardRouter(ctx *core.Context) {
	handler := dashboard.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "网站",
		Prefix: "/api/v1/site/dashboard",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/data",
			Handler:    handler.DashboardData,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func musicRouter(ctx *core.Context) {
	handler := music.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "音乐管理",
		Prefix: "/api/v1/site/music",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveMusic,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateMusic,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateMusicState,
			Middleware: nil,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteMusic,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:id",
			Handler:    handler.DetailMusic,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetMusicList,
			Middleware: nil,
		},
	}
	RegisterRouter(group, routes)
}
