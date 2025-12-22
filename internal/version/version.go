package version

// 这些变量会在编译时通过 -ldflags 注入
var (
	// Version 版本号
	Version = "dev"
	// GitCommit Git提交哈希
	GitCommit = "unknown"
	// BuildTime 构建时间
	BuildTime = "unknown"
	// GoVersion Go版本
	GoVersion = "unknown"
)

// Info 返回版本信息
func Info() map[string]string {
	return map[string]string{
		"version":   Version,
		"gitCommit": GitCommit,
		"buildTime": BuildTime,
		"goVersion": GoVersion,
	}
}

// String 返回格式化的版本字符串
func String() string {
	return Version + " (" + GitCommit[:min(7, len(GitCommit))] + ")"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
