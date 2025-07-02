package main

import (
	"context"
	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/internal/core"
	"gitee.com/jieepre/go-site/internal/middleware"
	pkg "gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/internal/router"
	"gitee.com/jieepre/go-site/internal/service"
	webTemplates "gitee.com/jieepre/go-site/internal/web"
	"gitee.com/jieepre/go-site/pkg/result"
	"gitee.com/jieepre/go-site/static/console"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	g       errgroup.Group
	Service = service.InitAppService()
	cfg     = &Config{
		AdminPort:    ":8000",
		ConsolePort:  ":8890",
		WebPort:      ":8589",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
)

type Config struct {
	AdminPort    string
	ConsolePort  string
	WebPort      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func main() {
	core.InitLog()
	adminServer := newHTTPServer(cfg.AdminPort, AdminRouter())
	consoleServer := newHTTPServer(cfg.ConsolePort, ConsoleRouter())
	webServer := newHTTPServer(cfg.WebPort, WebRouter())

	servers := []*http.Server{adminServer, consoleServer, webServer}

	setupSignalHandler(servers...)
	g.Go(func() error {
		return startServer(adminServer, "adminServer")
	})

	g.Go(func() error {
		return startServer(consoleServer, "consoleServer")
	})

	g.Go(func() error {
		return startServer(webServer, "webServer")
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func AdminRouter() http.Handler {
	// 1.创建路由
	engine := gin.New()
	//初始化配置
	engine.Use(middleware.Cors())
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	engine.NoRoute(func(c *gin.Context) {
		result.With(c, 404, "接口不存在", nil)
		return
	})
	engine.Use(gzip.Gzip(gzip.DefaultCompression)).Use(middleware.JwtVerify())
	cxt := &pkg.Context{
		Engine:  engine,
		Service: Service,
	}
	core.InitResource()
	core.Timer()
	router.RegisterAdminRouter(cxt)
	return engine
}

func WebRouter() http.Handler {
	tmpl := template.Must(template.New("").Funcs(pkg.TemplateFunc()).ParseFS(webTemplates.Fs, "**/*.html"))
	engine := gin.New()
	//初始化配置
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	engine.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	})
	engine.Use(middleware.Statistics())
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	cxt := &pkg.Context{
		Engine:  engine,
		Service: Service,
	}
	engine.SetHTMLTemplate(tmpl)
	router.RegisterWebRouter(cxt)
	return engine
}

func ConsoleRouter() http.Handler {
	engine := gin.New()
	//初始化配置
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	engine.StaticFS("/console", http.FS(console.Static))
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.POST("/console/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"baseUrl": config.Get().System.BaseUrl,
		})
	})
	return engine
}

func setupSignalHandler(servers ...*http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down servers...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		for _, server := range servers {
			if err := shutdownServer(server, shutdownCtx, "admin"); err != nil {
				log.Fatal("Server forced to shutdown:", err)
			}
		}

		log.Println("All servers exited")

	}()
}
func startServer(server *http.Server, name string) error {
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("%s server failed: %v", name, err)
	}
	log.Printf("%s server running on %s\n", name, server.Addr)
	return err
}

func shutdownServer(server *http.Server, ctx context.Context, name string) error {
	log.Printf("Shutting down %s server...\n", name)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down %s server: %v", name, err)
		return err
	}
	log.Printf("%s server gracefully stopped\n", name)
	return nil
}

func newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
