package app

import (
	"html/template"
	"net/http"
	"time"

	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/internal/middleware"
	"gitee.com/jieepre/go-site/internal/monitor"
	pkg "gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/internal/router"
	"gitee.com/jieepre/go-site/internal/service"
	webTemplates "gitee.com/jieepre/go-site/internal/web"
	"gitee.com/jieepre/go-site/pkg/result"
	"gitee.com/jieepre/go-site/static/console"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// CreateAdminRouter 创建管理后台路由
func CreateAdminRouter(service *service.AppService, healthChecker *monitor.HealthChecker, metrics *monitor.Metrics) http.Handler {
	engine := gin.New()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 添加中间件
	engine.Use(middleware.Cors())
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 错误处理中间件
	engine.Use(middleware.ErrorHandler())

	// 指标收集中间件
	engine.Use(monitor.MetricsMiddleware(metrics))

	// 限流中间件 - API限流
	engine.Use(middleware.APIRateLimit(100, time.Minute)) // 每分钟100次请求

	// 登录限流（更严格）
	engine.Use(middleware.LoginRateLimit())

	// 404处理
	engine.NoRoute(func(c *gin.Context) {
		result.With(c, http.StatusNotFound, "接口不存在", nil)
		return
	})

	// Gzip压缩
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	// JWT验证
	engine.Use(middleware.JwtVerify())

	// 监控端点（不需要JWT验证）
	engine.GET("/health", healthChecker.HealthCheck)
	engine.GET("/health/ready", healthChecker.ReadinessCheck)
	engine.GET("/health/live", healthChecker.LivenessCheck)
	engine.GET("/metrics", metrics.MetricsHandler)

	// 创建上下文
	cxt := &pkg.Context{
		Engine:  engine,
		Service: service,
	}

	// 注册路由
	router.RegisterAdminRouter(cxt)

	return engine
}

// CreateWebRouter 创建前台路由
func CreateWebRouter(service *service.AppService) http.Handler {
	// 解析模板
	tmpl := template.Must(template.New("").Funcs(pkg.TemplateFunc()).ParseFS(webTemplates.Fs, "**/*.html"))

	engine := gin.New()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 添加中间件
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 通用限流 - IP限流
	engine.Use(middleware.IPBasedRateLimit(60, time.Minute)) // 每分钟60次请求

	// 404处理
	engine.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	})

	// 统计中间件
	engine.Use(middleware.Statistics())

	// 缓存中间件（只缓存前台页面）
	engine.Use(middleware.CacheMiddleware(5 * time.Minute))

	// Gzip压缩
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	// 设置HTML模板
	engine.SetHTMLTemplate(tmpl)

	// 创建上下文
	cxt := &pkg.Context{
		Engine:  engine,
		Service: service,
	}

	// 注册路由
	router.RegisterWebRouter(cxt)

	return engine
}

// CreateConsoleRouter 创建控制台路由
func CreateConsoleRouter() http.Handler {
	engine := gin.New()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 添加中间件
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 静态文件服务
	engine.StaticFS("/console", http.FS(console.Static))

	// Gzip压缩
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	// 配置接口
	engine.POST("/console/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"baseUrl": config.Get().System.BaseUrl,
		})
	})

	return engine
}
