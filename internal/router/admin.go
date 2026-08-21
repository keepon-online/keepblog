package router

import (
	"net/http"

	"gitee.com/jieepre/keepblog/api/admin/about"
	"gitee.com/jieepre/keepblog/api/admin/category"
	"gitee.com/jieepre/keepblog/api/admin/common"
	"gitee.com/jieepre/keepblog/api/admin/dashboard"
	"gitee.com/jieepre/keepblog/api/admin/link"
	logviewer "gitee.com/jieepre/keepblog/api/admin/log"
	"gitee.com/jieepre/keepblog/api/admin/login"
	"gitee.com/jieepre/keepblog/api/admin/monitor"
	"gitee.com/jieepre/keepblog/api/admin/music"
	"gitee.com/jieepre/keepblog/api/admin/notice"
	"gitee.com/jieepre/keepblog/api/admin/post"
	"gitee.com/jieepre/keepblog/api/admin/system"
	"gitee.com/jieepre/keepblog/api/admin/tags"
	"gitee.com/jieepre/keepblog/api/admin/website"
	"gitee.com/jieepre/keepblog/internal/middleware"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	ws "gitee.com/jieepre/keepblog/internal/websocket"
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
	logViewerRouter(ctx)
	noticeRouter(ctx)
	websocketRouter(ctx)
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPost,
			Path:       "/api/change-password",
			Handler:    handler.ChangePassword,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/api/getInfo",
			Handler:    handler.GetInfo,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodPut,
			Path:       "/api/updateProfile",
			Handler:    handler.UpdateProfile,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/realtime",
			Handler:    handler.GetRealtime,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/general",
			Handler:    handler.General,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/loadavg",
			Handler:    handler.Loadavg,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/ram",
			Handler:    handler.RAM,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/cpu",
			Handler:    handler.CPU,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/net",
			Handler:    handler.Net,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/diskUsage",
			Handler:    handler.DiskUsage,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/diskIOStat",
			Handler:    handler.DiskIOStat,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/base/os",
			Handler:    handler.LoadDashboardOsInfo,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/base/:ioOption/:netOption",
			Handler:    handler.LoadDashboardBaseInfo,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/current/:ioOption/:netOption",
			Handler:    handler.LoadDashboardCurrentInfo,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateCategory,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateCategoryState,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:categoryId",
			Handler:    handler.DeleteCategory,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:categoryId",
			Handler:    handler.DetailCategory,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetCategoryList,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdatePost,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-publish",
			Handler:    handler.PublishPost,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-top",
			Handler:    handler.TopPost,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-cover/:postId",
			Handler:    handler.UpdatePostCoverImag,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodPut,
			Path:       "/update-cover",
			Handler:    handler.UpdatePostAllCoverImag,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:postId",
			Handler:    handler.DeletePost,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:postId",
			Handler:    handler.DetailPost,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetList,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
		// /list 供前台侧边栏 tag cloud 使用（无需登录）
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetTags,
			Middleware: []gin.HandlerFunc{},
		},
		// 以下为后台管理 CRUD，需登录
		{
			Method:     http.MethodGet,
			Path:       "/manage-list",
			Handler:    handler.GetTagList,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveTag,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateTag,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:tagId",
			Handler:    handler.DeleteTag,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:tagId",
			Handler:    handler.DetailTag,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateLink,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateLinkState,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteLink,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:id",
			Handler:    handler.DetailLink,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetLinkList,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		}, {
			Method:     http.MethodGet,
			Path:       "/access",
			Handler:    handler.AccessLog,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
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
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateMusic,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateMusicState,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteMusic,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:id",
			Handler:    handler.DetailMusic,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetMusicList,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
	}
	RegisterRouter(group, routes)
}

func logViewerRouter(ctx *core.Context) {
	handler := logviewer.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "日志管理",
		Prefix: "/api/log",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/level",
			Handler:    handler.GetLogLevel,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/level",
			Handler:    handler.SetLogLevel,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/stats",
			Handler:    handler.GetLogStats,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetLogList,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodGet,
			Path:       "/read",
			Handler:    handler.ReadLog,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
	}
	RegisterRouter(group, routes)
}

func noticeRouter(ctx *core.Context) {
	handler := notice.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "通知管理",
		Prefix: "/api/v1/notice",
	}
	group := ctx.Engine.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetNotices,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/read/:id",
			Handler:    handler.MarkAsRead,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPut,
			Path:       "/read-all",
			Handler:    handler.MarkAllAsRead,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteNotice,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/clear",
			Handler:    handler.ClearNotices,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
		{
			Method:     http.MethodPost,
			Path:       "/send",
			Handler:    handler.SendNotice,
			Middleware: []gin.HandlerFunc{middleware.JwtVerify()},
		},
	}
	RegisterRouter(group, routes)
}

func websocketRouter(ctx *core.Context) {
	ctx.Engine.GET("/ws", ws.HandleWebSocket)
}
