package main

import (
	"os"

	"gitee.com/jieepre/keepblog/internal/app"
	"github.com/gookit/slog"
)

func main() {
	// 创建应用实例
	application := app.New()

	// 初始化应用
	if err := application.Initialize(); err != nil {
		slog.Fatalf("Failed to initialize application: %v", err)
		// gookit/slog 的 Fatalf 只记日志不退出，必须显式退出
		os.Exit(1)
	}

	// 设置服务器
	application.SetupServers()

	// 运行应用
	if err := application.Run(); err != nil {
		slog.Fatalf("Application failed: %v", err)
		os.Exit(1)
	}
}
