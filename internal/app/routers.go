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

		// 注册前台路由
		router.RegisterWebRouter(webCtx)
	}

	// 后台管理路由组（使用特定中间件）
	adminGroup := engine.Group("/api")
	{
		// 后台专用中间件
		//adminGroup.Use(middleware.Cors())
		adminGroup.Use(middleware.ErrorHandler())
		adminGroup.Use(monitor.MetricsMiddleware(metrics))
		adminGroup.Use(middleware.APIRateLimit(100, time.Minute)) // 每分钟100次请求
		adminGroup.Use(middleware.LoginRateLimit())
		adminGroup.Use(middleware.JwtVerify())

		// 创建后台上下文
		adminCtx := &pkg.Context{
			Engine:  engine,
			Service: service,
		}

		// 注册后台路由
		router.RegisterAdminRouter(adminCtx)
	}

	return engine
}
