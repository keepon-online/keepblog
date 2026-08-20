package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/cache"
	appconfig "gitee.com/jieepre/go-site/internal/config"
	"gitee.com/jieepre/go-site/internal/core"
	"gitee.com/jieepre/go-site/internal/monitor"
	"gitee.com/jieepre/go-site/internal/pkg/oss"
	"gitee.com/jieepre/go-site/internal/service"
	"gitee.com/jieepre/go-site/internal/version"
	"gitee.com/jieepre/go-site/pkg/area"
	"gitee.com/jieepre/go-site/pkg/hash"
	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"github.com/gookit/slog"

	"github.com/go-co-op/gocron"
)

// Application 应用程序结构
type Application struct {
	config        *appconfig.ServerConfig
	service       *service.AppService
	servers       []*http.Server
	healthChecker *monitor.HealthChecker
	metrics       *monitor.Metrics
	scheduler     *gocron.Scheduler
}

// New 创建新的应用实例。配置相关字段在 Initialize（config.Load 之后）中填充。
func New() *Application {
	return &Application{
		service:       service.InitAppService(),
		servers:       make([]*http.Server, 0),
		healthChecker: monitor.NewHealthChecker(version.Version),
		metrics:       monitor.NewMetrics(),
	}
}

// Initialize 初始化应用。所有组件在此按依赖顺序显式初始化——
// 此前这些动作分散在各包的 init() 里（import 即连库/起协程/读配置），
// 既无法控制顺序，也让任何包都无法脱离运行环境被测试。
func (app *Application) Initialize() error {
	// 配置：先加载，后校验（jwt.secret 缺失等 fail-fast 在此拦截）
	if err := config.Load(); err != nil {
		return fmt.Errorf("load configuration failed: %w", err)
	}
	if err := config.ValidateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	cfg := config.Get()

	// 密钥类配置注入（hashids salt 留空即保持历史默认值）
	hash.Configure(cfg.Hashids.Salt)
	if err := jwttoken.Configure([]byte(cfg.Jwt.Secret)); err != nil {
		return fmt.Errorf("configure jwt failed: %w", err)
	}
	jwttoken.StartBlacklistCleanup()

	// 日志
	core.InitLog()

	// 可选组件：失败降级，不阻断启动
	area.EnsureIPDB()
	if err := oss.Init(); err != nil {
		slog.Errorf("Minio initialization failed: %v", err)
	}
	if err := cache.InitRedis(); err != nil {
		slog.Errorf("Redis initialization failed: %v", err)
	}

	// 数据库：建表/种子/索引（此前 core 包 init() 里 import 即连库）
	global.GORM = core.InitDB()
	core.InitResource()

	// 定时任务
	app.scheduler = core.Timer()

	// 配置就绪后才能确定监听地址
	app.config = appconfig.GetServerConfig()
	return nil
}

// SetupServers 设置服务器
func (app *Application) SetupServers() {
	server := app.newHTTPServer(app.config.AdminPort, app.createRouter())
	app.servers = []*http.Server{server}
}

// Run 运行应用，阻塞直到出错或收到退出信号。
// 退出时按序关闭 HTTP 服务、定时任务、Redis 连接与黑名单清理协程。
func (app *Application) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.startServer(app.servers[0], "admin")
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Println("Shutting down servers...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		for i, server := range app.servers {
			name := fmt.Sprintf("server-%d", i)
			if err := app.shutdownServer(server, shutdownCtx, name); err != nil {
				slog.Errorf("Error shutting down %s server: %v", name, err)
			}
		}
		if app.scheduler != nil {
			app.scheduler.Stop()
		}
		if err := cache.Close(); err != nil {
			slog.Warnf("Redis close failed: %v", err)
		}
		jwttoken.StopBlacklistCleanup()

		slog.Println("All servers exited")
		return nil
	}
}

// newHTTPServer 创建HTTP服务器
func (app *Application) newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  app.config.ReadTimeout,
		WriteTimeout: app.config.WriteTimeout,
	}
}

// startServer 启动服务器
func (app *Application) startServer(server *http.Server, name string) error {
	slog.Infof("%s server starting on %s\n", name, server.Addr)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Errorf("%s server failed: %v", name, err)
		return err
	}
	return nil
}

// shutdownServer 关闭服务器
func (app *Application) shutdownServer(server *http.Server, ctx context.Context, name string) error {
	slog.Printf("Shutting down %s server...", name)
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	slog.Printf("%s server gracefully stopped", name)
	return nil
}

func (app *Application) createRouter() http.Handler {
	return CreateRouters(app.service, app.healthChecker, app.metrics)
}
