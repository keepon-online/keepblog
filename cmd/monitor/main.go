package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gitee.com/jieepre/go-site/api/admin/monitor"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	// 创建Gin引擎
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// 添加CORS中间件
	engine.Use(corsMiddleware())

	// 添加日志中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 创建监控处理器
	ctx := &core.Context{
		Engine: engine,
	}
	handler := monitor.Handler{Context: ctx}

	// 静态文件服务
	engine.Static("/web", "./web")

	// 监控路由组
	monitorGroup := engine.Group("/api/monitor")
	{
		monitorGroup.GET("/server", handler.Monitor)
		monitorGroup.GET("/realtime", handler.GetRealtime)
		monitorGroup.GET("/general", handler.General)
		monitorGroup.GET("/cpu", handler.CPU)
		monitorGroup.GET("/ram", handler.RAM)
		monitorGroup.GET("/loadavg", handler.Loadavg)
		monitorGroup.GET("/net", handler.Net)
		monitorGroup.GET("/diskUsage", handler.DiskUsage)
		monitorGroup.GET("/diskIOStat", handler.DiskIOStat)
	}

	// 健康检查路由
	engine.GET("/health", func(c *gin.Context) {
		result.Ok(c, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"message":   "Monitor service is running",
		})
	})

	// 根路径重定向到监控页面
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/monitor.html")
	})

	// 启动服务器
	fmt.Println("🚀 监控服务启动成功!")
	fmt.Println("📊 监控面板: http://localhost:8090/web/monitor.html")
	fmt.Println("🔍 健康检查: http://localhost:8090/health")
	fmt.Println("📡 API文档:")
	fmt.Println("   - 完整监控: http://localhost:8090/api/monitor/server")
	fmt.Println("   - 实时数据: http://localhost:8090/api/monitor/realtime")
	fmt.Println("   - 系统信息: http://localhost:8090/api/monitor/general")
	fmt.Println("   - CPU信息: http://localhost:8090/api/monitor/cpu")
	fmt.Println("   - 内存信息: http://localhost:8090/api/monitor/ram")
	fmt.Println("   - 网络信息: http://localhost:8090/api/monitor/net")
	fmt.Println("   - 磁盘信息: http://localhost:8090/api/monitor/diskUsage")

	log.Fatal(engine.Run(":8090"))
}
