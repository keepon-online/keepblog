package app

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"gitee.com/jieepre/keepblog/internal/errors"
	"gitee.com/jieepre/keepblog/internal/middleware"
	"gitee.com/jieepre/keepblog/internal/monitor"
	pkg "gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/internal/router"
	"gitee.com/jieepre/keepblog/internal/service"
	webTemplates "gitee.com/jieepre/keepblog/internal/web"
	"gitee.com/jieepre/keepblog/static/console"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func CreateRouters(service *service.AppService, healthChecker *monitor.HealthChecker, metrics *monitor.Metrics) http.Handler {
	engine := gin.New()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 全局中间件
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	engine.Use(middleware.Cors())
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	// 静态文件服务 - 控制台
	engine.StaticFS("/console", http.FS(console.Static))

	// 404处理中间件 - 根据路径前缀决定返回哪种404响应
	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是API路径，返回JSON格式的404
		if strings.HasPrefix(path, "/api/") {
			_ = c.Error(errors.New(errors.ErrResourceNotFound, "接口不存在"))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// 否则返回HTML格式的404页面
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	// 健康检查和监控端点
	engine.GET("/health", healthChecker.HealthCheck)
	engine.GET("/health/ready", healthChecker.ReadinessCheck)
	engine.GET("/health/live", healthChecker.LivenessCheck)
	engine.GET("/metrics", metrics.MetricsHandler)

	// 前台路由组（使用特定中间件）
	// 注意：分组中间件只对"通过该分组注册"的路由生效。必须把 group 传给
	// RegisterWebRouter，由页面路由在其子分组上注册；此前路由各自在 Engine 上
	// 新建分组，导致本组的限流/统计/页面缓存从未生效过。
	webGroup := engine.Group("/")
	{
		// 前台专用中间件
		webGroup.Use(middleware.IPBasedRateLimit(60, time.Minute)) // 每分钟60次请求
		webGroup.Use(middleware.Statistics())
		webGroup.Use(middleware.CacheMiddleware(5 * time.Minute))

		// 设置HTML模板
		tmpl := template.Must(template.New("").Funcs(pkg.TemplateFunc()).ParseFS(webTemplates.Fs, "**/*.html"))
		engine.SetHTMLTemplate(tmpl)

		// 创建前台上下文
		webCtx := &pkg.Context{
			Engine:  engine,
			Service: service,
		}

		// 注册前台路由（页面路由挂在 webGroup 上继承中间件；
		// 静态资源/Feed/音乐 API 由 RegisterWebRouter 内部按各自缓存策略注册在 Engine 上）
		router.RegisterWebRouter(webCtx, webGroup)
	}

	// 后台管理路由组（使用特定中间件）
	// base 为空串以保持 admin.go 中的绝对前缀（/api/v1/...）不变。
	adminGroup := engine.Group("")
	{
		// 后台专用中间件
		adminGroup.Use(middleware.ErrorHandler())
		adminGroup.Use(monitor.MetricsMiddleware(metrics))
		adminGroup.Use(middleware.APIRateLimit(100, time.Minute)) // 每分钟100次请求
		adminGroup.Use(middleware.LoginRateLimit())
		// 鉴权在分组级统一生效（JwtVerify 自带 login/refreshToken 等白名单），
		// 路由级不再重复挂载
		adminGroup.Use(middleware.JwtVerify())

		// 创建后台上下文
		adminCtx := &pkg.Context{
			Engine:  engine,
			Service: service,
		}

		// 注册后台路由
		router.RegisterAdminRouter(adminCtx, adminGroup)
	}

	return engine
}
