package monitor

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psNet "github.com/shirou/gopsutil/v3/net"

	"gitee.com/jieepre/keepblog/pkg/result"
)

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
	// 获取各核心使用率（200ms采样兼顾精确度与接口毫秒级高响应）
	corePercents, err := cpu.Percent(200*time.Millisecond, true)
	if err != nil {
		corePercents = []float64{}
	}

	var usage float64
	if len(corePercents) > 0 {
		var sum float64
		for _, p := range corePercents {
			sum += p
		}
		usage = sum / float64(len(corePercents))
	} else {
		if percents, err := cpu.Percent(0, false); err == nil && len(percents) > 0 {
			usage = percents[0]
		}
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

	// 获取Swap信息
	swapStat, _ := mem.SwapMemory()

	memInfo := &MemoryInfo{
		Total:       memStat.Total,
		Used:        memStat.Used,
		Free:        memStat.Free,
		Available:   memStat.Available,
		UsedPercent: h.formatPercent(memStat.UsedPercent),
		TotalFormat: h.formatBytes(memStat.Total),
		UsedFormat:  h.formatBytes(memStat.Used),
		FreeFormat:  h.formatBytes(memStat.Free),
	}

	// 添加Swap信息
	if swapStat != nil {
		memInfo.SwapTotal = swapStat.Total
		memInfo.SwapUsed = swapStat.Used
		memInfo.SwapFree = swapStat.Free
		memInfo.SwapUsedPercent = h.formatPercent(swapStat.UsedPercent)
		memInfo.SwapTotalFormat = h.formatBytes(swapStat.Total)
		memInfo.SwapUsedFormat = h.formatBytes(swapStat.Used)
	}

	return memInfo, nil
}

// 获取磁盘信息
func (h *Handler) getDiskInfo() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	// 过滤虚拟文件系统、loop 镜像设备与临时挂载
	var candidates []disk.PartitionStat
	for _, partition := range partitions {
		// 过滤 loop 虚拟块设备（如 snap 软件包挂载的沙箱镜像）
		if strings.HasPrefix(partition.Device, "/dev/loop") {
			continue
		}
		if h.isVirtualFilesystem(partition.Fstype) || h.isTempPartition(partition.Mountpoint) {
			continue
		}
		// 过滤 Docker 内部 overlay 衍生路径
		if strings.Contains(partition.Mountpoint, "docker/overlay2") {
			continue
		}
		candidates = append(candidates, partition)
	}

	// 设备初筛：同一设备多个挂载点保留优先级最高的挂载点
	deviceMap := make(map[string]disk.PartitionStat)
	for _, partition := range candidates {
		if existing, ok := deviceMap[partition.Device]; ok {
			if mountPriority(partition.Mountpoint) < mountPriority(existing.Mountpoint) {
				deviceMap[partition.Device] = partition
			}
		} else {
			deviceMap[partition.Device] = partition
		}
	}

	// 智能去重：仅当容量、已用空间与 Inode 总数完全重合时，才判定为容器/bind mount 的同一物理磁盘映射
	// 避免将同等规格的多块独立物理磁盘误判去重
	dedupMap := make(map[string]DiskInfo)
	for _, partition := range deviceMap {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		// 过滤无容量或极小无意义微型分区（< 100MB）
		if usage.Total < 100*1024*1024 {
			continue
		}

		info := DiskInfo{
			Device:            partition.Device,
			Mountpoint:        partition.Mountpoint,
			Fstype:            partition.Fstype,
			Total:             usage.Total,
			Used:              usage.Used,
			Free:              usage.Free,
			UsedPercent:       h.formatPercent(usage.UsedPercent),
			TotalFormat:       h.formatBytes(usage.Total),
			UsedFormat:        h.formatBytes(usage.Used),
			FreeFormat:        h.formatBytes(usage.Free),
			InodesTotal:       usage.InodesTotal,
			InodesUsed:        usage.InodesUsed,
			InodesFree:        usage.InodesFree,
			InodesUsedPercent: h.formatPercent(usage.InodesUsedPercent),
		}

		// 复合特征签名：容量 + 已用 + Inode 总数完全重合才是同一个底层文件系统
		signature := fmt.Sprintf("%d_%d_%d", usage.Total, usage.Used, usage.InodesTotal)
		if existing, ok := dedupMap[signature]; ok {
			if mountPriority(info.Mountpoint) < mountPriority(existing.Mountpoint) {
				dedupMap[signature] = info
			}
		} else {
			dedupMap[signature] = info
		}
	}

	var diskInfos []DiskInfo
	for _, info := range dedupMap {
		diskInfos = append(diskInfos, info)
	}

	sort.Slice(diskInfos, func(i, j int) bool {
		return diskInfos[i].Mountpoint < diskInfos[j].Mountpoint
	})

	return diskInfos, nil
}

// mountPriority 挂载点展示优先级，越小越优先（/ 优于 /app，/app 优于其子目录）
func mountPriority(mountpoint string) int {
	switch {
	case mountpoint == "/":
		return 0
	case mountpoint == "/app" || mountpoint == "/data" || mountpoint == "/home":
		return 1
	case strings.HasPrefix(mountpoint, "/var/"):
		return 3
	case strings.HasPrefix(mountpoint, "/app/"):
		return 5
	default:
		return 10
	}
}

// isVirtualFilesystem 判断是否为虚拟文件系统
// 注意 overlay 不在此列：容器根目录的 fstype 就是 overlay，需要保留，交给挂载点/总量去重处理
func (h *Handler) isVirtualFilesystem(fstype string) bool {
	virtualFilesystems := []string{
		"tmpfs", "devtmpfs", "sysfs", "proc", "cgroup",
		"cgroup2", "pstore", "squashfs",
		"debugfs", "tracefs", "securityfs", "sockfs",
		"pipefs", "rpc_pipefs", "rpc_pipe", "binfmt_misc",
		"devpts", "ramfs", "hugetlbfs", "mqueue", "autofs",
		"iso9660", "udf", "nsfs", "bpf", "configfs", "fusectl",
	}

	for _, vfs := range virtualFilesystems {
		if fstype == vfs {
			return true
		}
	}
	return false
}

// isTempPartition 判断是否为临时分区或无关紧要的分区
func (h *Handler) isTempPartition(mountpoint string) bool {
	tempPartitions := []string{
		"/dev", "/run", "/sys", "/proc", "/tmp", "/var/tmp",
		"/boot/efi", "/snap", "/var/lib/docker", "/var/lib/lxc",
		"/var/lib/kubelet", "/var/lib/containerd", "/lost+found",
		"/etc", // Docker 注入容器的 bind mount（resolv.conf/hostname/hosts）
	}

	for _, tp := range tempPartitions {
		if mountpoint == tp || (len(mountpoint) > len(tp) && strings.HasPrefix(mountpoint, tp+"/")) {
			return true
		}
	}
	return false
}

// 获取网络信息
func (h *Handler) getNetworkInfo() ([]NetworkInfo, error) {
	interfaces, err := psNet.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var networkInfos []NetworkInfo
	for _, iface := range interfaces {
		// 过滤掉虚拟网络接口和无用接口
		if h.isVirtualInterface(iface.Name) {
			continue
		}
		// 检查网络接口状态
		isUp := h.isNetworkInterfaceUp(iface.Name)
		// 如果接口未启用且没有流量，则跳过
		if !isUp && iface.BytesRecv == 0 && iface.BytesSent == 0 {
			continue
		}

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

// isVirtualInterface 判断是否为虚拟网络接口
func (h *Handler) isVirtualInterface(name string) bool {
	virtualInterfaces := []string{
		"lo", "docker", "br-", "veth", "tun", "tap", "virbr",
		"vmnet", "vboxnet", "ppp", "ip6gre", "ipip", "sit",
		"gre", "stf", "gif", "dummy", "nlmon", "zt",
	}

	for _, prefix := range virtualInterfaces {
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			return true
		}
	}
	return false
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
	totalProcs := 0
	if hostInfo, err := host.Info(); err == nil && hostInfo.Procs > 0 {
		totalProcs = int(hostInfo.Procs)
	}
	if totalProcs == 0 {
		totalProcs = runtime.NumGoroutine()
	}
	return &ProcessInfo{
		Total:    totalProcs,
		Running:  runtime.NumGoroutine(),
		Sleeping: 0,
		Stopped:  0,
		Zombie:   0,
	}, nil
}

// 获取实时统计数据
func (h *Handler) getRealtimeStats() (*RealtimeStats, error) {
	// CPU使用率（200ms采样）
	cpuPercents, err := cpu.Percent(200*time.Millisecond, false)
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

	// 磁盘使用率(取去重后的代表分区；直接取 partitions[0] 在容器里可能是 Docker 注入的 bind mount)
	var diskPercent float64
	if diskInfos, err := h.getDiskInfo(); err == nil && len(diskInfos) > 0 {
		diskPercent = diskInfos[0].UsedPercent
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
func (h *Handler) formatPercent(percent float64) float64 {
	sprintf := fmt.Sprintf("%.1f", percent)
	float, _ := strconv.ParseFloat(sprintf, 64)
	return float
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
