package monitor

import "gitee.com/jieepre/keepblog/internal/pkg/core"

// Handler 监控接口处理器
type Handler struct {
	*core.Context
}

// SystemInfo 系统信息统一结构
type SystemInfo struct {
	General   GeneralInfo   `json:"general"`
	CPU       CPUInfo       `json:"cpu"`
	Memory    MemoryInfo    `json:"memory"`
	Disk      []DiskInfo    `json:"disk"`
	Network   []NetworkInfo `json:"network"`
	Load      LoadInfo      `json:"load"`
	Process   ProcessInfo   `json:"process"`
	Timestamp int64         `json:"timestamp"`
}

// GeneralInfo 系统基本信息
type GeneralInfo struct {
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	Kernel       string `json:"kernel"`
	Uptime       int64  `json:"uptime"`
	UptimeFormat string `json:"uptimeFormat"`
	CPUCores     int    `json:"cpuCores"`
	GoVersion    string `json:"goVersion"`
}

// CPUInfo CPU信息
type CPUInfo struct {
	UsagePercent float64   `json:"usagePercent"`
	Cores        int       `json:"cores"`
	ModelName    string    `json:"modelName"`
	Frequency    float64   `json:"frequency"`
	CoreDetails  []float64 `json:"coreDetails"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total           uint64  `json:"total"`
	Used            uint64  `json:"used"`
	Free            uint64  `json:"free"`
	Available       uint64  `json:"available"`
	UsedPercent     float64 `json:"usedPercent"`
	TotalFormat     string  `json:"totalFormat"`
	UsedFormat      string  `json:"usedFormat"`
	FreeFormat      string  `json:"freeFormat"`
	SwapTotal       uint64  `json:"swapTotal"`
	SwapUsed        uint64  `json:"swapUsed"`
	SwapFree        uint64  `json:"swapFree"`
	SwapUsedPercent float64 `json:"swapUsedPercent"`
	SwapTotalFormat string  `json:"swapTotalFormat"`
	SwapUsedFormat  string  `json:"swapUsedFormat"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
	TotalFormat string  `json:"totalFormat"`
	UsedFormat  string  `json:"usedFormat"`
	FreeFormat  string  `json:"freeFormat"`
}

// NetworkInfo 网络接口信息
type NetworkInfo struct {
	Name        string `json:"name"`
	BytesRecv   uint64 `json:"bytesRecv"`
	BytesSent   uint64 `json:"bytesSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	RecvFormat  string `json:"recvFormat"`
	SentFormat  string `json:"sentFormat"`
	IsUp        bool   `json:"isUp"`
}

// LoadInfo 系统负载信息
type LoadInfo struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	Total    int `json:"total"`
	Running  int `json:"running"`
	Sleeping int `json:"sleeping"`
	Stopped  int `json:"stopped"`
	Zombie   int `json:"zombie"`
}

// RealtimeStats 实时统计数据
type RealtimeStats struct {
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
	DiskPercent   float64 `json:"diskPercent"`
	NetworkIn     uint64  `json:"networkIn"`
	NetworkOut    uint64  `json:"networkOut"`
	Timestamp     int64   `json:"timestamp"`
}
