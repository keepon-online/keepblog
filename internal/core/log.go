package core

import (
	"fmt"

	"gitee.com/jieepre/go-site/internal/logger"
	"github.com/gookit/slog"
)

// InitLog 初始化结构化日志（同时输出到控制台与 ./logs 文件，自动轮转）
func InitLog() {
	logConfig := logger.LogConfig{
		Level:      "info",
		Format:     "text",
		Output:     "both", // 同时输出到控制台和文件
		Path:       "./logs",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	if err := logger.InitLogger(logConfig); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		// 退回到默认日志配置
		slog.Configure(func(logger *slog.SugaredLogger) {
			f := logger.Formatter.(*slog.TextFormatter)
			f.EnableColor = true
		})
	}
}
