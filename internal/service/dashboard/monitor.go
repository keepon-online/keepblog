package dashboard

import (
	"encoding/json"
	"sort"
	"strings"
	"sync"
	"time"

	"gitee.com/jieepre/go-site/internal/model/response"
	"gitee.com/jieepre/go-site/pkg/cmd"
	"gitee.com/jieepre/go-site/pkg/copier"
	"gitee.com/jieepre/go-site/pkg/xpack"
	"github.com/gookit/slog"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func (s Service) LoadOsInfo() (*response.OsInfo, error) {
	var baseInfo response.OsInfo
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	baseInfo.OS = hostInfo.OS
	baseInfo.Platform = hostInfo.Platform
	baseInfo.PlatformFamily = hostInfo.PlatformFamily
	baseInfo.KernelArch = hostInfo.KernelArch
	baseInfo.KernelVersion = hostInfo.KernelVersion

	if baseInfo.KernelArch == "armv7l" {
		baseInfo.KernelArch = "armv7"
	}
	if baseInfo.KernelArch == "x86_64" {
		baseInfo.KernelArch = "amd64"
	}
	return &baseInfo, nil
}

func (u Service) LoadBaseInfo(ioOption string, netOption string) (*response.DashboardBase, error) {
	var baseInfo response.DashboardBase
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	baseInfo.Hostname = hostInfo.Hostname
	baseInfo.OS = hostInfo.OS
	baseInfo.Platform = hostInfo.Platform
	baseInfo.PlatformFamily = hostInfo.PlatformFamily
	baseInfo.PlatformVersion = hostInfo.PlatformVersion
	baseInfo.KernelArch = hostInfo.KernelArch
	baseInfo.KernelVersion = hostInfo.KernelVersion
	ss, _ := json.Marshal(hostInfo)
	baseInfo.VirtualizationSystem = string(ss)

	cpuInfo, err := cpu.Info()
	if err == nil {
		baseInfo.CPUModelName = cpuInfo[0].ModelName
	}

	baseInfo.CPUCores, _ = cpu.Counts(false)
	baseInfo.CPULogicalCores, _ = cpu.Counts(true)

	baseInfo.CurrentInfo = *u.LoadCurrentInfo(ioOption, netOption)
	return &baseInfo, nil
}

func (u *Service) LoadCurrentInfo(ioOption string, netOption string) *response.DashboardCurrent {
	var currentInfo response.DashboardCurrent
	hostInfo, _ := host.Info()
	currentInfo.Uptime = hostInfo.Uptime
	currentInfo.TimeSinceUptime = time.Now().Add(-time.Duration(hostInfo.Uptime) * time.Second).Format("2006-01-02 15:04:05")
	currentInfo.Procs = hostInfo.Procs

	currentInfo.CPUTotal, _ = cpu.Counts(true)
	totalPercent, _ := cpu.Percent(0, false)
	if len(totalPercent) == 1 {
		currentInfo.CPUUsedPercent = totalPercent[0]
		currentInfo.CPUUsed = currentInfo.CPUUsedPercent * 0.01 * float64(currentInfo.CPUTotal)
	}
	currentInfo.CPUPercent, _ = cpu.Percent(0, true)

	loadInfo, _ := load.Avg()
	currentInfo.Load1 = loadInfo.Load1
	currentInfo.Load5 = loadInfo.Load5
	currentInfo.Load15 = loadInfo.Load15
	currentInfo.LoadUsagePercent = loadInfo.Load1 / (float64(currentInfo.CPUTotal*2) * 0.75) * 100

	memoryInfo, _ := mem.VirtualMemory()
	currentInfo.MemoryTotal = memoryInfo.Total
	currentInfo.MemoryAvailable = memoryInfo.Available
	currentInfo.MemoryUsed = memoryInfo.Used
	currentInfo.MemoryUsedPercent = memoryInfo.UsedPercent

	swapInfo, _ := mem.SwapMemory()
	currentInfo.SwapMemoryTotal = swapInfo.Total
	currentInfo.SwapMemoryAvailable = swapInfo.Free
	currentInfo.SwapMemoryUsed = swapInfo.Used
	currentInfo.SwapMemoryUsedPercent = swapInfo.UsedPercent

	currentInfo.DiskData = loadDiskInfo()
	currentInfo.GPUData = loadGPUInfo()

	if ioOption == "all" {
		diskInfo, _ := disk.IOCounters()
		for _, state := range diskInfo {
			currentInfo.IOReadBytes += state.ReadBytes
			currentInfo.IOWriteBytes += state.WriteBytes
			currentInfo.IOCount += (state.ReadCount + state.WriteCount)
			currentInfo.IOReadTime += state.ReadTime
			currentInfo.IOWriteTime += state.WriteTime
		}
	} else {
		diskInfo, _ := disk.IOCounters(ioOption)
		for _, state := range diskInfo {
			currentInfo.IOReadBytes += state.ReadBytes
			currentInfo.IOWriteBytes += state.WriteBytes
			currentInfo.IOCount += (state.ReadCount + state.WriteCount)
			currentInfo.IOReadTime += state.ReadTime
			currentInfo.IOWriteTime += state.WriteTime
		}
	}

	if netOption == "all" {
		netInfo, _ := net.IOCounters(false)
		if len(netInfo) != 0 {
			currentInfo.NetBytesSent = netInfo[0].BytesSent
			currentInfo.NetBytesRecv = netInfo[0].BytesRecv
		}
	} else {
		netInfo, _ := net.IOCounters(true)
		for _, state := range netInfo {
			if state.Name == netOption {
				currentInfo.NetBytesSent = state.BytesSent
				currentInfo.NetBytesRecv = state.BytesRecv
			}
		}
	}

	currentInfo.ShotTime = time.Now()
	return &currentInfo
}

type diskInfo struct {
	Type   string
	Mount  string
	Device string
}

func loadDiskInfo() []response.DiskInfo {
	var datas []response.DiskInfo

	// 优先使用 gopsutil 的 Partitions 获取分区信息（跨平台且更可靠）
	partitions, err := disk.Partitions(false) // false = 只获取物理分区
	if err == nil && len(partitions) > 0 {
		return loadDiskInfoFromPartitions(partitions)
	}

	// 备选方案：使用 df 命令
	stdout, err := cmd.ExecWithTimeOut("df -hT -P|grep '/'|grep -v tmpfs|grep -v 'snap/core'|grep -v udev", 2*time.Second)
	if err != nil {
		stdout, err = cmd.ExecWithTimeOut("df -lhT -P|grep '/'|grep -v tmpfs|grep -v 'snap/core'|grep -v udev", 1*time.Second)
		if err != nil {
			return datas
		}
	}
	lines := strings.Split(stdout, "\n")

	var mounts []diskInfo
	// Docker 容器内常见的虚拟挂载点，应排除
	var excludes = []string{
		"/mnt/cdrom", "/boot", "/boot/efi", "/dev", "/dev/shm",
		"/run/lock", "/run", "/run/shm", "/run/user",
		"/etc/hostname", "/etc/hosts", "/etc/resolv.conf", // Docker 虚拟挂载
		"/proc", "/sys",
	}
	// Docker 容器内应排除的挂载点前缀
	var excludePrefixes = []string{
		"/proc/", "/sys/", "/dev/", "/etc/",
		"/app/data", "/app/logs", "/app/config", // Docker volume 挂载
	}

	// 用于按设备去重：device -> 最佳挂载点信息
	deviceMap := make(map[string]diskInfo)
	// 挂载点优先级（越小越优先）
	mountPriority := func(mount string) int {
		switch {
		case mount == "/":
			return 0
		case mount == "/app" || mount == "/data" || mount == "/home":
			return 1
		case strings.HasPrefix(mount, "/var/"):
			return 3
		case strings.HasPrefix(mount, "/app/") || strings.HasPrefix(mount, "/data/") || strings.HasPrefix(mount, "/home/"):
			return 50 // Docker volume 子目录，较低优先级
		case strings.HasPrefix(mount, "/etc/"):
			return 100 // Docker 虚拟挂载，最低优先级
		default:
			return 10
		}
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		if fields[1] == "tmpfs" || fields[1] == "overlay" || fields[1] == "shm" {
			continue
		}
		if strings.Contains(fields[2], "M") || strings.Contains(fields[2], "K") {
			continue
		}
		if strings.Contains(fields[6], "docker") {
			continue
		}

		mount := fields[6]

		// 检查排除列表
		isExclude := false
		for _, exclude := range excludes {
			if exclude == mount {
				isExclude = true
				break
			}
		}
		if isExclude {
			continue
		}

		// 检查排除前缀
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(mount, prefix) {
				isExclude = true
				break
			}
		}
		if isExclude {
			continue
		}

		device := fields[0]
		fsType := fields[1]
		newMount := diskInfo{Type: fsType, Device: device, Mount: mount}

		// 按设备去重：保留优先级更高的挂载点
		if existing, ok := deviceMap[device]; ok {
			if mountPriority(mount) < mountPriority(existing.Mount) {
				deviceMap[device] = newMount
			}
		} else {
			deviceMap[device] = newMount
		}
	}

	// 将去重后的设备转为 slice
	for _, info := range deviceMap {
		mounts = append(mounts, info)
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(len(mounts))
	for i := 0; i < len(mounts); i++ {
		go func(timeoutCh <-chan time.Time, mount diskInfo) {
			defer wg.Done()

			var itemData response.DiskInfo
			itemData.Path = mount.Mount
			itemData.Type = mount.Type
			itemData.Device = mount.Device
			select {
			case <-timeoutCh:
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
				slog.Errorf("load disk info from %s failed, err: timeout", mount.Mount)
			default:
				state, err := disk.Usage(mount.Mount)
				if err != nil {
					mu.Lock()
					datas = append(datas, itemData)
					mu.Unlock()
					slog.Errorf("load disk info from %s failed, err: %v", mount.Mount, err)
					return
				}
				itemData.Total = state.Total
				itemData.Free = state.Free
				itemData.Used = state.Used
				itemData.UsedPercent = state.UsedPercent
				itemData.InodesTotal = state.InodesTotal
				itemData.InodesUsed = state.InodesUsed
				itemData.InodesFree = state.InodesFree
				itemData.InodesUsedPercent = state.InodesUsedPercent
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
			}
		}(time.After(5*time.Second), mounts[i])
	}
	wg.Wait()

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Path < datas[j].Path
	})
	return datas
}

// loadDiskInfoFromPartitions 使用 gopsutil 分区信息（更可靠，自动去重）
func loadDiskInfoFromPartitions(partitions []disk.PartitionStat) []response.DiskInfo {
	var datas []response.DiskInfo

	// 按设备去重
	deviceMap := make(map[string]disk.PartitionStat)

	// 排除的文件系统类型
	excludeFsTypes := map[string]bool{
		"tmpfs": true, "devtmpfs": true, "overlay": true, "shm": true,
		"squashfs": true, "iso9660": true, "udf": true,
	}

	// 排除的挂载点
	excludeMounts := map[string]bool{
		"/dev": true, "/dev/shm": true, "/proc": true, "/sys": true,
		"/run": true, "/run/lock": true, "/boot": true, "/boot/efi": true,
		"/etc/hostname": true, "/etc/hosts": true, "/etc/resolv.conf": true,
	}

	// 排除的挂载点前缀
	excludePrefixes := []string{"/proc/", "/sys/", "/dev/", "/run/", "/snap/"}

	// 挂载点优先级
	mountPriority := func(mount string) int {
		switch {
		case mount == "/":
			return 0
		case mount == "/app" || mount == "/data" || mount == "/home":
			return 1
		case strings.HasPrefix(mount, "/var/"):
			return 3
		default:
			return 10
		}
	}

	for _, p := range partitions {
		// 排除特定文件系统
		if excludeFsTypes[p.Fstype] {
			continue
		}

		// 排除特定挂载点
		if excludeMounts[p.Mountpoint] {
			continue
		}

		// 排除前缀匹配
		isExclude := false
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(p.Mountpoint, prefix) {
				isExclude = true
				break
			}
		}
		if isExclude {
			continue
		}

		// 按设备去重
		if existing, ok := deviceMap[p.Device]; ok {
			if mountPriority(p.Mountpoint) < mountPriority(existing.Mountpoint) {
				deviceMap[p.Device] = p
			}
		} else {
			deviceMap[p.Device] = p
		}
	}

	// 获取磁盘使用情况
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	for _, p := range deviceMap {
		wg.Add(1)
		go func(partition disk.PartitionStat) {
			defer wg.Done()

			var itemData response.DiskInfo
			itemData.Path = partition.Mountpoint
			itemData.Type = partition.Fstype
			itemData.Device = partition.Device

			state, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				slog.Errorf("load disk info from %s failed, err: %v", partition.Mountpoint, err)
				return
			}

			// 过滤太小的分区（< 1GB）
			if state.Total < 1024*1024*1024 {
				return
			}

			itemData.Total = state.Total
			itemData.Free = state.Free
			itemData.Used = state.Used
			itemData.UsedPercent = state.UsedPercent
			itemData.InodesTotal = state.InodesTotal
			itemData.InodesUsed = state.InodesUsed
			itemData.InodesFree = state.InodesFree
			itemData.InodesUsedPercent = state.InodesUsedPercent

			mu.Lock()
			datas = append(datas, itemData)
			mu.Unlock()
		}(p)
	}
	wg.Wait()

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Path < datas[j].Path
	})
	return datas
}

func loadGPUInfo() []response.GPUInfo {
	list := xpack.LoadGpuInfo()
	if len(list) == 0 {
		return nil
	}
	var data []response.GPUInfo
	for _, gpu := range list {
		var dataItem response.GPUInfo
		if err := copier.Copy(&dataItem, &gpu); err != nil {
			continue
		}
		dataItem.PowerUsage = dataItem.PowerDraw + " / " + dataItem.MaxPowerLimit
		dataItem.MemoryUsage = dataItem.MemUsed + " / " + dataItem.MemTotal
		data = append(data, dataItem)
	}
	return data
}
