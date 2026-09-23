package router

import (
	"net/http"

	"gitee.com/jieepre/keepblog/api/admin/about"
	ai "gitee.com/jieepre/keepblog/api/admin/ai"
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
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	ws "gitee.com/jieepre/keepblog/internal/websocket"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(ctx *core.Context, group *gin.RouterGroup) {
	userRouter(ctx, group)
	monitorRouter(ctx, group)
	categoryRouter(ctx, group)
	tagAdminRouter(ctx, group)
	postRouter(ctx, group)
	commonRouter(ctx, group)
	aboutBackendRouter(ctx, group)
	linkBackendRouter(ctx, group)
	websiteBackendRouter(ctx, group)
	logsRouter(ctx, group)
	dashboardRouter(ctx, group)
	musicRouter(ctx, group)
	logViewerRouter(ctx, group)
	noticeRouter(ctx, group)
	aiRouter(ctx, group)
	websocketRouter(ctx, group)
}

func userRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := login.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "系统管理",
		Prefix: "/",
	}
	group = group.Group(routeGroup.Prefix)
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
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPost,
			Path:       "/api/change-password",
			Handler:    handler.ChangePassword,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
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
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/api/getInfo",
			Handler:    handler.GetInfo,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodPut,
			Path:       "/api/updateProfile",
			Handler:    handler.UpdateProfile,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func monitorRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := monitor.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "监控",
		Prefix: "/api/monitor",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/server",
			Handler:    handler.Monitor,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/realtime",
			Handler:    handler.GetRealtime,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/general",
			Handler:    handler.General,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/loadavg",
			Handler:    handler.Loadavg,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/ram",
			Handler:    handler.RAM,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/cpu",
			Handler:    handler.CPU,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/net",
			Handler:    handler.Net,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/diskUsage",
			Handler:    handler.DiskUsage,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/diskIOStat",
			Handler:    handler.DiskIOStat,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/base/os",
			Handler:    handler.LoadDashboardOsInfo,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/base/:ioOption/:netOption",
			Handler:    handler.LoadDashboardBaseInfo,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/current/:ioOption/:netOption",
			Handler:    handler.LoadDashboardCurrentInfo,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func categoryRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := category.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "分类管理",
		Prefix: "/api/v1/site/category",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveCategory,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateCategory,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateCategoryState,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:categoryId",
			Handler:    handler.DeleteCategory,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:categoryId",
			Handler:    handler.DetailCategory,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetCategoryList,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func postRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := post.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "测试分组",
		Prefix: "/api/v1/site/post",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SavePost,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdatePost,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-publish",
			Handler:    handler.PublishPost,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-top",
			Handler:    handler.TopPost,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-cover/:postId",
			Handler:    handler.UpdatePostCoverImag,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodPut,
			Path:       "/update-cover",
			Handler:    handler.UpdatePostAllCoverImag,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:postId",
			Handler:    handler.DeletePost,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:postId",
			Handler:    handler.DetailPost,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetList,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/random-cover",
			Handler:    handler.GetRandomCover,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func tagAdminRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := tags.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "标签管理",
		Prefix: "/api/v1/site/tags",
	}
	group = group.Group(routeGroup.Prefix)
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
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveTag,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateTag,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:tagId",
			Handler:    handler.DeleteTag,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:tagId",
			Handler:    handler.DetailTag,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func commonRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := common.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "公用",
		Prefix: "/api/v1",
	}
	group = group.Group(routeGroup.Prefix)
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

func aboutBackendRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := about.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "关于",
		Prefix: "/api/v1/site/about",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.Save,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/detail",
			Handler:    handler.Detail,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func linkBackendRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := link.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "友链管理",
		Prefix: "/api/v1/site/link",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveLink,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateLink,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateLinkState,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteLink,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:id",
			Handler:    handler.DetailLink,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetLinkList,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func websiteBackendRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := website.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "网站",
		Prefix: "/api/v1/site/web",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.SaveWebSite,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/detail",
			Handler:    handler.GetWebSite,
			Middleware: []gin.HandlerFunc{},
		},
	}
	RegisterRouter(group, routes)
}

func logsRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := system.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "网站",
		Prefix: "/api/v1/site/logs",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/logon",
			Handler:    handler.LoginLog,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		}, {
			Method:     http.MethodGet,
			Path:       "/access",
			Handler:    handler.AccessLog,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}
func dashboardRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := dashboard.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "网站",
		Prefix: "/api/v1/site/dashboard",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/data",
			Handler:    handler.DashboardData,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func musicRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := music.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "音乐管理",
		Prefix: "/api/v1/site/music",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodPost,
			Path:       "/save",
			Handler:    handler.SaveMusic,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update",
			Handler:    handler.UpdateMusic,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/update-state",
			Handler:    handler.UpdateMusicState,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteMusic,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/detail/:id",
			Handler:    handler.DetailMusic,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetMusicList,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func logViewerRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := logviewer.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "日志管理",
		Prefix: "/api/log",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/level",
			Handler:    handler.GetLogLevel,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/level",
			Handler:    handler.SetLogLevel,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/stats",
			Handler:    handler.GetLogStats,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetLogList,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/read",
			Handler:    handler.ReadLog,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func aiRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := ai.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "AI 辅助写作",
		Prefix: "/api/v1/ai",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/status",
			Handler:    handler.Status,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPost,
			Path:       "/edit",
			Handler:    handler.Edit,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodGet,
			Path:       "/templates",
			Handler:    handler.Templates,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/templates",
			Handler:    handler.SaveTemplate,
			Middleware: nil,
		},
		{
			Method:     http.MethodGet,
			Path:       "/config",
			Handler:    handler.AIConfig,
			Middleware: nil,
		},
		{
			Method:     http.MethodPut,
			Path:       "/config",
			Handler:    handler.SaveAIConfig,
			Middleware: nil,
		},
	}
	RegisterRouter(group, routes)
}

func noticeRouter(ctx *core.Context, group *gin.RouterGroup) {
	handler := notice.Handler{Context: ctx}
	routeGroup := RouteGroup{
		Name:   "通知管理",
		Prefix: "/api/v1/notice",
	}
	group = group.Group(routeGroup.Prefix)
	routes := []Route{
		{
			Method:     http.MethodGet,
			Path:       "/list",
			Handler:    handler.GetNotices,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/read/:id",
			Handler:    handler.MarkAsRead,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPut,
			Path:       "/read-all",
			Handler:    handler.MarkAllAsRead,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodDelete,
			Path:       "/delete/:id",
			Handler:    handler.DeleteNotice,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodDelete,
			Path:       "/clear",
			Handler:    handler.ClearNotices,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
		{
			Method:     http.MethodPost,
			Path:       "/send",
			Handler:    handler.SendNotice,
			Middleware: nil, // 鉴权由 adminGroup 分组级统一生效
		},
	}
	RegisterRouter(group, routes)
}

func websocketRouter(ctx *core.Context, group *gin.RouterGroup) {
	ctx.Engine.GET("/ws", ws.HandleWebSocket)
}
