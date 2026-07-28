package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

// LogConfig 日志配置
type LogConfig struct {
	Level       string `yaml:"level" default:"info"`
	Format      string `yaml:"format" default:"text"` // json, text
	Output      string `yaml:"output" default:"both"` // file, console, both
	Path        string `yaml:"path" default:"./logs"`
	MaxSize     int    `yaml:"maxSize" default:"100"`   // MB
	MaxBackups  int    `yaml:"maxBackups" default:"10"` // 保留文件数
	MaxAge      int    `yaml:"maxAge" default:"30"`     // 天数
	Compress    bool   `yaml:"compress" default:"true"`
	AsyncBuffer int    `yaml:"asyncBuffer" default:"1000"` // 异步缓冲区大小
	ShowCaller  bool   `yaml:"showCaller" default:"true"`  // 显示调用者信息
}

// 全局变量
var (
	currentConfig  LogConfig
	configMutex    sync.RWMutex
	asyncLogChan   chan logEntry
	asyncWaitGroup sync.WaitGroup
	asyncRunning   bool
	asyncMutex     sync.Mutex
)

// logEntry 异步日志条目
type logEntry struct {
	Level   slog.Level
	Message string
	Fields  map[string]any
	Time    time.Time
	Caller  string
}

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	configMutex.Lock()
	currentConfig = config
	configMutex.Unlock()

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
			f.TimeFormat = time.DateTime
		}
	})

	// 设置日志级别
	level := parseLogLevel(config.Level)
	slog.SetLogLevel(level)

	// 配置输出处理器
	if config.Output == "console" {
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
		slog.PushHandler(fileHandler)
		slog.PushHandler(errorHandler)
	} else {
		slog.Reset()
		slog.PushHandler(fileHandler)
		slog.PushHandler(errorHandler)
	}

	// 启动异步写入
	if config.AsyncBuffer > 0 {
		startAsyncWriter(config.AsyncBuffer)
	}

	// 启动日志清理
	go startLogCleanup(config)

	return nil
}

// startAsyncWriter 启动异步写入器
func startAsyncWriter(bufferSize int) {
	asyncMutex.Lock()
	defer asyncMutex.Unlock()

	if asyncRunning {
		return
	}

	asyncLogChan = make(chan logEntry, bufferSize)
	asyncRunning = true
	asyncWaitGroup.Add(1)

	go func() {
		defer asyncWaitGroup.Done()
		for entry := range asyncLogChan {
			fields := entry.Fields
			if fields == nil {
				fields = make(map[string]any)
			}
			if entry.Caller != "" {
				fields["caller"] = entry.Caller
			}

			switch entry.Level {
			case slog.DebugLevel:
				slog.WithFields(fields).Debug(entry.Message)
			case slog.InfoLevel:
				slog.WithFields(fields).Info(entry.Message)
			case slog.WarnLevel:
				slog.WithFields(fields).Warn(entry.Message)
			case slog.ErrorLevel:
				slog.WithFields(fields).Error(entry.Message)
			default:
				slog.WithFields(fields).Info(entry.Message)
			}
		}
	}()
}

// StopAsyncWriter 停止异步写入器
func StopAsyncWriter() {
	asyncMutex.Lock()
	defer asyncMutex.Unlock()

	if !asyncRunning {
		return
	}

	close(asyncLogChan)
	asyncWaitGroup.Wait()
	asyncRunning = false
}

// AsyncLog 异步日志写入
func AsyncLog(level slog.Level, msg string, fields map[string]any) {
	if !asyncRunning {
		// 如果异步未启动，直接写入
		slog.WithFields(fields).Log(level, msg)
		return
	}

	caller := ""
	configMutex.RLock()
	showCaller := currentConfig.ShowCaller
	configMutex.RUnlock()

	if showCaller {
		caller = getCaller(3)
	}

	entry := logEntry{
		Level:   level,
		Message: msg,
		Fields:  fields,
		Time:    time.Now(),
		Caller:  caller,
	}

	select {
	case asyncLogChan <- entry:
	default:
		// 缓冲区满，直接写入
		slog.WithFields(fields).Log(level, msg)
	}
}

// getCaller 获取调用者信息
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	// 只取文件名
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' || file[i] == '\\' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

// GetLevel 获取当前日志级别
func GetLevel() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return currentConfig.Level
}

// SetLevel 动态设置日志级别
func SetLevel(level string) error {
	parsedLevel := parseLogLevel(level)
	slog.SetLogLevel(parsedLevel)

	configMutex.Lock()
	currentConfig.Level = level
	configMutex.Unlock()

	slog.Infof("Log level changed to: %s", level)
	return nil
}

// GetConfig 获取当前日志配置
func GetConfig() LogConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return currentConfig
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
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

// startLogCleanup 启动日志清理
func startLogCleanup(config LogConfig) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// 启动时先执行一次清理
	cleanupLogs(config)

	for range ticker.C {
		cleanupLogs(config)
	}
}

// cleanupLogs 清理过期日志
func cleanupLogs(config LogConfig) {
	logDir := config.Path
	maxAge := time.Duration(config.MaxAge) * 24 * time.Hour
	maxBackups := config.MaxBackups

	files, err := os.ReadDir(logDir)
	if err != nil {
		slog.Errorf("Failed to read log directory: %v", err)
		return
	}

	var logFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			logFiles = append(logFiles, file)
		}
	}

	// 按修改时间排序
	sort.Slice(logFiles, func(i, j int) bool {
		infoI, _ := logFiles[i].Info()
		infoJ, _ := logFiles[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	now := time.Now()
	deleted := 0

	for i, file := range logFiles {
		info, err := file.Info()
		if err != nil {
			continue
		}

		filePath := filepath.Join(logDir, file.Name())

		// 删除超过保留数量的文件
		if i >= maxBackups {
			if err := os.Remove(filePath); err == nil {
				deleted++
				slog.Debugf("Deleted old log file: %s", file.Name())
			}
			continue
		}

		// 删除超过保留天数的文件
		if now.Sub(info.ModTime()) > maxAge {
			if err := os.Remove(filePath); err == nil {
				deleted++
				slog.Debugf("Deleted expired log file: %s", file.Name())
			}
		}
	}

	if deleted > 0 {
		slog.Infof("Log cleanup completed, deleted %d files", deleted)
	}
}

// LogStats 日志统计信息
type LogStats struct {
	TotalFiles   int      `json:"totalFiles"`
	TotalSize    int64    `json:"totalSize"`
	TotalSizeStr string   `json:"totalSizeStr"`
	OldestFile   string   `json:"oldestFile"`
	NewestFile   string   `json:"newestFile"`
	LogFiles     []string `json:"logFiles"`
}

// GetLogStats 获取日志统计信息
func GetLogStats() (*LogStats, error) {
	configMutex.RLock()
	logDir := currentConfig.Path
	configMutex.RUnlock()

	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	stats := &LogStats{
		LogFiles: make([]string, 0),
	}

	var oldestTime, newestTime time.Time
	var oldestName, newestName string

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		stats.TotalFiles++
		stats.TotalSize += info.Size()
		stats.LogFiles = append(stats.LogFiles, file.Name())

		if oldestTime.IsZero() || info.ModTime().Before(oldestTime) {
			oldestTime = info.ModTime()
			oldestName = file.Name()
		}
		if newestTime.IsZero() || info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newestName = file.Name()
		}
	}

	stats.TotalSizeStr = formatBytes(stats.TotalSize)
	stats.OldestFile = oldestName
	stats.NewestFile = newestName

	return stats, nil
}

// ReadLogFile 读取日志文件内容
func ReadLogFile(filename string, lines int, search string) ([]string, error) {
	configMutex.RLock()
	logDir := currentConfig.Path
	configMutex.RUnlock()

	filePath := filepath.Join(logDir, filename)

	// 安全检查，防止路径遍历
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(logDir)) {
		return nil, fmt.Errorf("invalid file path")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var result []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if search == "" || strings.Contains(line, search) {
			result = append(result, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	// 返回最后N行
	if lines > 0 && len(result) > lines {
		result = result[len(result)-lines:]
	}

	return result, nil
}

// formatBytes 格式化字节数
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
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
	AsyncLog(slog.InfoLevel, msg, l.buildFields(fields...))
}

// Warn 记录警告日志
func (l *StructuredLogger) Warn(msg string, fields ...any) {
	AsyncLog(slog.WarnLevel, msg, l.buildFields(fields...))
}

// Error 记录错误日志
func (l *StructuredLogger) Error(msg string, fields ...any) {
	AsyncLog(slog.ErrorLevel, msg, l.buildFields(fields...))
}

// Debug 记录调试日志
func (l *StructuredLogger) Debug(msg string, fields ...any) {
	AsyncLog(slog.DebugLevel, msg, l.buildFields(fields...))
}

// buildFields 构建日志字段
func (l *StructuredLogger) buildFields(fields ...any) map[string]any {
	result := map[string]any{
		"component": l.component,
	}

	// 处理键值对
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok && i+1 < len(fields) {
			result[key] = fields[i+1]
		}
	}

	return result
}
