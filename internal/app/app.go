package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/internal/cache"
	appconfig "gitee.com/jieepre/go-site/internal/config"
	"gitee.com/jieepre/go-site/internal/core"
	"gitee.com/jieepre/go-site/internal/monitor"
	"gitee.com/jieepre/go-site/internal/service"
	"github.com/gookit/slog"

	"golang.org/x/sync/errgroup"
)

// Application 应用程序结构
type Application struct {
	config        *appconfig.ServerConfig
	service       *service.AppService
	servers       []*http.Server
	healthChecker *monitor.HealthChecker
	metrics       *monitor.Metrics
}

// New 创建新的应用实例
func New() *Application {
	return &Application{
		config:        appconfig.GetServerConfig(),
		service:       service.InitAppService(),
		servers:       make([]*http.Server, 0),
		healthChecker: monitor.NewHealthChecker("1.0.0"), // 版本号
		metrics:       monitor.NewMetrics(),
	}
}

// Initialize 初始化应用
func (app *Application) Initialize() error {
	// 验证配置
	if err := config.ValidateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// 初始化日志
	core.InitLog()

	// 初始化Redis
	if err := cache.InitRedis(); err != nil {
		slog.Errorf("Redis initialization failed: %v", err)
	}

	// 初始化资源
	core.InitResource()

	// 启动定时任务
	core.Timer()

	return nil
}

// SetupServers 设置服务器
func (app *Application) SetupServers() {
	server := app.newHTTPServer(app.config.AdminPort, app.createRouter())
	app.servers = []*http.Server{server}
}

// Run 运行应用
func (app *Application) Run() error {
	var g errgroup.Group

	// 设置优雅关闭
	app.setupGracefulShutdown()

	if err := app.startServer(app.servers[0], "admin"); err != nil {
		slog.Errorf("Failed to start %s server: %v", "admin", err)
		return err
	}

	return g.Wait()
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

// setupGracefulShutdown 设置优雅关闭
func (app *Application) setupGracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		slog.Println("Shutting down servers...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		for i, server := range app.servers {
			serverNames := []string{"admin", "console", "web"}
			name := fmt.Sprintf("server-%d", i) // 默认名称
			if i < len(serverNames) {
				name = serverNames[i]
			}
			if err := app.shutdownServer(server, shutdownCtx, name); err != nil {
				slog.Printf("Error shutting down %s server: %v", name, err)
			}
		}

		slog.Println("All servers exited")
		os.Exit(0)
	}()
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
