package monitor

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psNet "github.com/shirou/gopsutil/v3/net"
)

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
	Uptime       int64  `json:"uptime"`       // 秒
	UptimeFormat string `json:"uptimeFormat"` // 格式化显示
	CPUCores     int    `json:"cpuCores"`
	GoVersion    string `json:"goVersion"`
}

// CPUInfo CPU信息
type CPUInfo struct {
	UsagePercent float64   `json:"usagePercent"` // 使用率百分比
	Cores        int       `json:"cores"`        // 核心数
	ModelName    string    `json:"modelName"`    // CPU型号
	Frequency    float64   `json:"frequency"`    // 频率 MHz
	CoreDetails  []float64 `json:"coreDetails"`  // 各核心使用率
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total       uint64  `json:"total"`       // 总内存 bytes
	Used        uint64  `json:"used"`        // 已使用 bytes
	Free        uint64  `json:"free"`        // 空闲内存 bytes
	Available   uint64  `json:"available"`   // 可用内存 bytes
	UsedPercent float64 `json:"usedPercent"` // 使用率百分比
	TotalFormat string  `json:"totalFormat"` // 格式化显示
	UsedFormat  string  `json:"usedFormat"`  // 格式化显示
	FreeFormat  string  `json:"freeFormat"`  // 格式化显示
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Device      string  `json:"device"`      // 设备名
	Mountpoint  string  `json:"mountpoint"`  // 挂载点
	Fstype      string  `json:"fstype"`      // 文件系统类型
	Total       uint64  `json:"total"`       // 总容量 bytes
	Used        uint64  `json:"used"`        // 已使用 bytes
	Free        uint64  `json:"free"`        // 空闲容量 bytes
	UsedPercent float64 `json:"usedPercent"` // 使用率百分比
	TotalFormat string  `json:"totalFormat"` // 格式化显示
	UsedFormat  string  `json:"usedFormat"`  // 格式化显示
	FreeFormat  string  `json:"freeFormat"`  // 格式化显示
}

// NetworkInfo 网络接口信息
type NetworkInfo struct {
	Name        string `json:"name"`        // 接口名
	BytesRecv   uint64 `json:"bytesRecv"`   // 接收字节数
	BytesSent   uint64 `json:"bytesSent"`   // 发送字节数
	PacketsRecv uint64 `json:"packetsRecv"` // 接收包数
	PacketsSent uint64 `json:"packetsSent"` // 发送包数
	RecvFormat  string `json:"recvFormat"`  // 格式化显示
	SentFormat  string `json:"sentFormat"`  // 格式化显示
	IsUp        bool   `json:"isUp"`        // 接口状态
}

// LoadInfo 系统负载信息
type LoadInfo struct {
	Load1  float64 `json:"load1"`  // 1分钟负载
	Load5  float64 `json:"load5"`  // 5分钟负载
	Load15 float64 `json:"load15"` // 15分钟负载
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	Total    int `json:"total"`    // 总进程数
	Running  int `json:"running"`  // 运行中进程数
	Sleeping int `json:"sleeping"` // 睡眠进程数
	Stopped  int `json:"stopped"`  // 停止进程数
	Zombie   int `json:"zombie"`   // 僵尸进程数
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

// Monitor 获取完整系统监控信息
func (h *Handler) Monitor(c *gin.Context) {
	systemInfo, err := h.getSystemInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取系统信息失败: %v", err))
		return
	}
	result.Ok(c, systemInfo)
}

// GetRealtime 获取实时监控数据
func (h *Handler) GetRealtime(c *gin.Context) {
	stats, err := h.getRealtimeStats()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取实时数据失败: %v", err))
		return
	}
	result.Ok(c, stats)
}

// General 获取系统基本信息
func (h *Handler) General(c *gin.Context) {
	generalInfo, err := h.getGeneralInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取系统信息失败: %v", err))
		return
	}
	result.Ok(c, generalInfo)
}

// CPU 获取CPU信息
func (h *Handler) CPU(c *gin.Context) {
	cpuInfo, err := h.getCPUInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取CPU信息失败: %v", err))
		return
	}
	result.Ok(c, cpuInfo)
}

// RAM 获取内存信息
func (h *Handler) RAM(c *gin.Context) {
	memInfo, err := h.getMemoryInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取内存信息失败: %v", err))
		return
	}
	result.Ok(c, memInfo)
}

// Loadavg 获取系统负载信息
func (h *Handler) Loadavg(c *gin.Context) {
	loadInfo, err := h.getLoadInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取系统负载信息失败: %v", err))
		return
	}
	result.Ok(c, loadInfo)
}

// Net 获取网络信息
func (h *Handler) Net(c *gin.Context) {
	netInfo, err := h.getNetworkInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取网络信息失败: %v", err))
		return
	}
	result.Ok(c, netInfo)
}

// DiskUsage 获取磁盘使用信息
func (h *Handler) DiskUsage(c *gin.Context) {
	diskInfo, err := h.getDiskInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取磁盘信息失败: %v", err))
		return
	}
	result.Ok(c, diskInfo)
}

// DiskIOStat 磁盘IO统计(兼容性接口)
func (h *Handler) DiskIOStat(c *gin.Context) {
	// 为了兼容旧接口，返回磁盘信息
	diskInfo, err := h.getDiskInfo()
	if err != nil {
		result.Error(c, fmt.Sprintf("获取磁盘信息失败: %v", err))
		return
	}

	// 转换为旧格式
	header := []string{"Device", "Total", "Used", "Free", "Use%", "Mounted"}
	var list [][]string
	for _, disk := range diskInfo {
		row := []string{
			disk.Device,
			disk.TotalFormat,
			disk.UsedFormat,
			disk.FreeFormat,
			fmt.Sprintf("%.1f%%", disk.UsedPercent),
			disk.Mountpoint,
		}
		list = append(list, row)
	}

	result.Ok(c, struct {
		Header []string   `json:"header"`
		List   [][]string `json:"list"`
	}{
		Header: header,
		List:   list,
	})
}

// 实现方法：获取完整系统信息
func (h *Handler) getSystemInfo() (*SystemInfo, error) {
	generalInfo, err := h.getGeneralInfo()
	if err != nil {
		return nil, fmt.Errorf("general info: %w", err)
	}

	cpuInfo, err := h.getCPUInfo()
	if err != nil {
		return nil, fmt.Errorf("cpu info: %w", err)
	}

	memInfo, err := h.getMemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("memory info: %w", err)
	}

	diskInfo, err := h.getDiskInfo()
	if err != nil {
		return nil, fmt.Errorf("disk info: %w", err)
	}

	netInfo, err := h.getNetworkInfo()
	if err != nil {
		return nil, fmt.Errorf("network info: %w", err)
	}

	loadInfo, err := h.getLoadInfo()
	if err != nil {
		return nil, fmt.Errorf("load info: %w", err)
	}

	processInfo, err := h.getProcessInfo()
	if err != nil {
		return nil, fmt.Errorf("process info: %w", err)
	}

	return &SystemInfo{
		General:   *generalInfo,
		CPU:       *cpuInfo,
		Memory:    *memInfo,
		Disk:      diskInfo,
		Network:   netInfo,
		Load:      *loadInfo,
		Process:   *processInfo,
		Timestamp: time.Now().Unix(),
	}, nil
}

// 获取系统基本信息
func (h *Handler) getGeneralInfo() (*GeneralInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = hostInfo.Hostname
	}

	return &GeneralInfo{
		Hostname:     hostname,
		OS:           fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion),
		Arch:         runtime.GOARCH,
		Kernel:       hostInfo.KernelVersion,
		Uptime:       int64(hostInfo.Uptime),
		UptimeFormat: h.formatUptime(int64(hostInfo.Uptime)),
		CPUCores:     runtime.NumCPU(),
		GoVersion:    runtime.Version(),
	}, nil
}

// 获取CPU信息
func (h *Handler) getCPUInfo() (*CPUInfo, error) {
	// 获取CPU使用率
	percents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// 获取各核心使用率
	corePercents, err := cpu.Percent(time.Second, true)
	if err != nil {
		corePercents = []float64{}
	}

	// 获取CPU信息
	infos, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var modelName string
	var frequency float64
	if len(infos) > 0 {
		modelName = infos[0].ModelName
		frequency = infos[0].Mhz
	}

	var usage float64
	if len(percents) > 0 {
		usage = percents[0]
	}

	return &CPUInfo{
		UsagePercent: usage,
		Cores:        runtime.NumCPU(),
		ModelName:    modelName,
		Frequency:    frequency,
		CoreDetails:  corePercents,
	}, nil
}

// 获取内存信息
func (h *Handler) getMemoryInfo() (*MemoryInfo, error) {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryInfo{
		Total:       memStat.Total,
		Used:        memStat.Used,
		Free:        memStat.Free,
		Available:   memStat.Available,
		UsedPercent: memStat.UsedPercent,
		TotalFormat: h.formatBytes(memStat.Total),
		UsedFormat:  h.formatBytes(memStat.Used),
		FreeFormat:  h.formatBytes(memStat.Free),
	}, nil
}

// 获取磁盘信息
func (h *Handler) getDiskInfo() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskInfos []DiskInfo
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		diskInfos = append(diskInfos, DiskInfo{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Fstype:      partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
			TotalFormat: h.formatBytes(usage.Total),
			UsedFormat:  h.formatBytes(usage.Used),
			FreeFormat:  h.formatBytes(usage.Free),
		})
	}

	return diskInfos, nil
}

// 获取网络信息
func (h *Handler) getNetworkInfo() ([]NetworkInfo, error) {
	interfaces, err := psNet.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var networkInfos []NetworkInfo
	for _, iface := range interfaces {
		// 检查网络接口状态
		isUp := h.isNetworkInterfaceUp(iface.Name)

		networkInfos = append(networkInfos, NetworkInfo{
			Name:        iface.Name,
			BytesRecv:   iface.BytesRecv,
			BytesSent:   iface.BytesSent,
			PacketsRecv: iface.PacketsRecv,
			PacketsSent: iface.PacketsSent,
			RecvFormat:  h.formatBytes(iface.BytesRecv),
			SentFormat:  h.formatBytes(iface.BytesSent),
			IsUp:        isUp,
		})
	}

	return networkInfos, nil
}

// 获取系统负载信息
func (h *Handler) getLoadInfo() (*LoadInfo, error) {
	loadStat, err := load.Avg()
	if err != nil {
		// Windows不支持负载信息，返回默认值
		return &LoadInfo{
			Load1:  0.0,
			Load5:  0.0,
			Load15: 0.0,
		}, nil
	}

	return &LoadInfo{
		Load1:  loadStat.Load1,
		Load5:  loadStat.Load5,
		Load15: loadStat.Load15,
	}, nil
}

// 获取进程信息
func (h *Handler) getProcessInfo() (*ProcessInfo, error) {
	// 简化实现，只返回基本信息
	return &ProcessInfo{
		Total:    runtime.NumGoroutine(), // 使用goroutine数量作为示例
		Running:  1,
		Sleeping: 0,
		Stopped:  0,
		Zombie:   0,
	}, nil
}

// 获取实时统计数据
func (h *Handler) getRealtimeStats() (*RealtimeStats, error) {
	// CPU使用率
	cpuPercents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	var cpuPercent float64
	if len(cpuPercents) > 0 {
		cpuPercent = cpuPercents[0]
	}

	// 内存使用率
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 磁盘使用率(取第一个磁盘)
	var diskPercent float64
	partitions, _ := disk.Partitions(false)
	if len(partitions) > 0 {
		usage, err := disk.Usage(partitions[0].Mountpoint)
		if err == nil {
			diskPercent = usage.UsedPercent
		}
	}

	// 网络流量
	netStats, _ := psNet.IOCounters(false)
	var networkIn, networkOut uint64
	if len(netStats) > 0 {
		networkIn = netStats[0].BytesRecv
		networkOut = netStats[0].BytesSent
	}

	return &RealtimeStats{
		CPUPercent:    cpuPercent,
		MemoryPercent: memStat.UsedPercent,
		DiskPercent:   diskPercent,
		NetworkIn:     networkIn,
		NetworkOut:    networkOut,
		Timestamp:     time.Now().Unix(),
	}, nil
}

// 工具函数：格式化字节数
func (h *Handler) formatBytes(bytes uint64) string {
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

// 工具函数：格式化运行时间
func (h *Handler) formatUptime(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	mins := (seconds % 3600) / 60
	secs := seconds % 60

	if days > 0 {
		return fmt.Sprintf("%d天 %d小时 %d分钟", days, hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时 %d分钟 %d秒", hours, mins, secs)
	} else {
		return fmt.Sprintf("%d分钟 %d秒", mins, secs)
	}
}

// 工具函数：检查网络接口状态
func (h *Handler) isNetworkInterfaceUp(name string) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, iface := range interfaces {
		if iface.Name == name {
			return iface.Flags&net.FlagUp != 0
		}
	}
	return false
}

// Dashboard相关的兼容性接口（保留）
func (h *Handler) LoadDashboardOsInfo(c *gin.Context) {
	generalInfo, err := h.getGeneralInfo()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, generalInfo)
}

func (h *Handler) LoadDashboardBaseInfo(c *gin.Context) {
	systemInfo, err := h.getSystemInfo()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, systemInfo)
}

func (h *Handler) LoadDashboardCurrentInfo(c *gin.Context) {
	stats, err := h.getRealtimeStats()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, stats)
}
