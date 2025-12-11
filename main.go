package main

import (
	"gitee.com/jieepre/go-site/internal/app"
	"github.com/gookit/slog"
)

func main() {
	// 创建应用实例
	application := app.New()

	// 初始化应用
	if err := application.Initialize(); err != nil {
		slog.Fatalf("Failed to initialize application: %v", err)
	}

	// 设置服务器
	application.SetupServers()

	// 运行应用
	if err := application.Run(); err != nil {
		slog.Fatalf("Application failed: %v", err)
	}
}
