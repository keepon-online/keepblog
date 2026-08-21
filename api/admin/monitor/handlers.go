package monitor

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gitee.com/jieepre/go-site/pkg/result"
)

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
	for _, info := range diskInfo {
		row := []string{
			info.Device,
			info.TotalFormat,
			info.UsedFormat,
			info.FreeFormat,
			fmt.Sprintf("%.1f%%", info.UsedPercent),
			info.Mountpoint,
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
