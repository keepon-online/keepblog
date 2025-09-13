package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level" default:"info"`
	Format     string `yaml:"format" default:"json"` // json, text
	Output     string `yaml:"output" default:"file"` // file, console, both
	Path       string `yaml:"path" default:"./logs"`
	MaxSize    int    `yaml:"maxSize" default:"100"`   // MB
	MaxBackups int    `yaml:"maxBackups" default:"10"` // 保留文件数
	MaxAge     int    `yaml:"maxAge" default:"30"`     // 天数
	Compress   bool   `yaml:"compress" default:"true"`
}

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	// 创建日志目录
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 配置日志格式
	slog.Configure(func(logger *slog.SugaredLogger) {
		if config.Format == "json" {
			logger.Formatter = slog.NewJSONFormatter()
		} else {
			f := logger.Formatter.(*slog.TextFormatter)
			f.EnableColor = config.Output == "console" || config.Output == "both"
			f.TimeFormat = "2006-01-02 15:04:05"
		}
	})

	// 设置日志级别
	level := parseLogLevel(config.Level)
	slog.SetLogLevel(level)

	// 配置输出处理器
	if config.Output == "console" {
		// 只输出到控制台
		return nil
	}

	// 文件输出配置
	logFilePath := filepath.Join(config.Path, "app.log")

	// 按日期轮转的文件处理器
	fileHandler, err := handler.NewRotateFileHandler(
		logFilePath,
		rotatefile.EveryDay,
		handler.WithLogLevels(slog.AllLevels),
	)
	if err != nil {
		return fmt.Errorf("failed to create file handler: %w", err)
	}

	// 错误日志单独文件
	errorLogPath := filepath.Join(config.Path, "error.log")
	errorHandler, err := handler.NewRotateFileHandler(
		errorLogPath,
		rotatefile.EveryDay,
		handler.WithLogLevels([]slog.Level{slog.ErrorLevel, slog.FatalLevel, slog.PanicLevel}),
	)
	if err != nil {
		return fmt.Errorf("failed to create error handler: %w", err)
	}

	if config.Output == "both" {
		// 同时输出到控制台和文件
		slog.PushHandler(fileHandler)
		slog.PushHandler(errorHandler)
	} else {
		// 只输出到文件
		slog.Reset()
		slog.PushHandler(fileHandler)
		slog.PushHandler(errorHandler)
	}

	return nil
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.DebugLevel
	case "info":
		return slog.InfoLevel
	case "warn", "warning":
		return slog.WarnLevel
	case "error":
		return slog.ErrorLevel
	case "fatal":
		return slog.FatalLevel
	case "panic":
		return slog.PanicLevel
	default:
		return slog.InfoLevel
	}
}

// StructuredLogger 结构化日志记录器
type StructuredLogger struct {
	component string
}

// NewStructuredLogger 创建结构化日志记录器
func NewStructuredLogger(component string) *StructuredLogger {
	return &StructuredLogger{component: component}
}

// Info 记录信息日志
func (l *StructuredLogger) Info(msg string, fields ...any) {
	slog.WithFields(l.buildFields(fields...)).Info(msg)
}

// Warn 记录警告日志
func (l *StructuredLogger) Warn(msg string, fields ...any) {
	slog.WithFields(l.buildFields(fields...)).Warn(msg)
}

// Error 记录错误日志
func (l *StructuredLogger) Error(msg string, fields ...any) {
	slog.WithFields(l.buildFields(fields...)).Error(msg)
}

// Debug 记录调试日志
func (l *StructuredLogger) Debug(msg string, fields ...any) {
	slog.WithFields(l.buildFields(fields...)).Debug(msg)
}

// buildFields 构建日志字段
func (l *StructuredLogger) buildFields(fields ...any) map[string]any {
	result := map[string]any{
		"component": l.component,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// 处理键值对
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok && i+1 < len(fields) {
			result[key] = fields[i+1]
		}
	}

	return result
}
