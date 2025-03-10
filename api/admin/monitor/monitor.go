package monitor

import (
	"bytes"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
	"os/exec"
	"strconv"
	"strings"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Monitor(c *gin.Context) {
	s := pkg.NewServer()
	result.Ok(c, &s)
}

// General info
// uname -mrs
// cat /etc/os-release
// nproc --all
// hostname
// uptime -p
func (h *Handler) General(c *gin.Context) {
	type GeneralInfo struct {
		KernelVersion string `json:"kernelVersion"`
		OS            string `json:"os"`
		Cores         string `json:"cores"`
		Hostname      string `json:"hostname"`
		Uptime        string `json:"uptime"`
	}

	var generalInfo GeneralInfo
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("uname", "-rs")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.KernelVersion = "Unknown"
	} else {
		generalInfo.KernelVersion = ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("bash", "-c", "cat /etc/os-release | grep \"PRETTY_NAME\" | awk -F\\\" '{print $2}'")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.OS = "Unknown"
	} else {
		generalInfo.OS = ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("nproc", "--all")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		println(err.Error())
		generalInfo.Cores = "Unknown"
	} else {
		generalInfo.Cores = ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("hostname")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.Hostname = "Unknown"
	} else {
		generalInfo.Hostname = ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("uptime", "-p")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.Uptime = "Unknown"
	} else {
		generalInfo.Uptime = ClearNewline(stdout.String())
	}
	result.Ok(c, generalInfo)
}

// Loadavg info
// cat /proc/loadavg
func (h *Handler) Loadavg(c *gin.Context) {
	type LoadavgInfo struct {
		Avg   string `json:"avg"`
		Avg5  string `json:"avg5"`
		Avg15 string `json:"avg15"`
		Cores string `json:"cores"`
	}

	var loadavgInfo LoadavgInfo
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/loadavg")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		procLoadavg := strings.Split(stdout.String(), " ")
		loadavgInfo.Avg = procLoadavg[0]
		loadavgInfo.Avg5 = procLoadavg[1]
		loadavgInfo.Avg15 = procLoadavg[2]
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("nproc", "--all")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		println(err.Error())
		loadavgInfo.Cores = ""
	} else {
		loadavgInfo.Cores = ClearNewline(stdout.String())
	}

	result.Ok(c, loadavgInfo)
}

// RAM info
// cat /proc/meminfo
func (h *Handler) RAM(c *gin.Context) {
	type RAMInfo struct {
		Total int `json:"total"`
		Free  int `json:"free"`
	}

	var ramInfo RAMInfo
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", "head -n 2 /proc/meminfo | awk -F \" \" '{print $2}'")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		procMeminfo := strings.Split(ClearNewline(stdout.String()), "\n")
		ramInfo.Total, err = strconv.Atoi(procMeminfo[0])
		if err == nil {
			ramInfo.Total *= 1000
		}

		ramInfo.Free, err = strconv.Atoi(procMeminfo[1])
		if err == nil {
			ramInfo.Free *= 1000
		}
	}

	result.Ok(c, ramInfo)
}

// CPU info
// cat /proc/stat
func (h *Handler) CPU(c *gin.Context) {
	var cpuList [][]string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/stat")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {

	} else {
		for _, line := range strings.Split(ClearNewline(stdout.String()), "\n") {
			if strings.Contains(line, "cpu") {
				cpuList = append(cpuList, strings.Fields(line))
			}
		}
	}
	result.Ok(c, cpuList)
}

// Net info
// cat /proc/net/dev
func (h *Handler) Net(c *gin.Context) {
	var netList [][]string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/net/dev")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		for _, line := range strings.Split(ClearNewline(stdout.String()), "\n")[2:] {
			fields := strings.Fields(line)
			netList = append(netList, fields)
		}

	}

	result.Ok(c, netList)
}

// DiskUsage info
// df -h  --output=source,size,used,avail,pcent,target,itotal,iused,iavail,ipcent,fstype
func (h *Handler) DiskUsage(c *gin.Context) {
	// Filesystem            Size  Used Avail Use% Mounted on Inodes IUsed IFree IUse% Type
	type DiskUsageInfo struct {
		Filesystem string `json:"filesystem"`
		Size       string `json:"size"`
		Used       string `json:"used"`
		Avail      string `json:"avail"`
		UsedPcent  string `json:"usedPcent"`
		MountedOn  string `json:"mountedOn"`
		Inodes     string `json:"inodes"`
		IUsed      string `json:"iUsed"`
		IFree      string `json:"iFree"`
		IUsedPcent string `json:"iUsedPcent"`
		Type       string `json:"type"`
	}
	var diskUsageList []DiskUsageInfo

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("df", "-h", "--output=size,used,avail,pcent,target,itotal,iused,iavail,ipcent,fstype,source")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {

	} else {
		for _, line := range strings.Split(ClearNewline(stdout.String()), "\n")[1:] {
			field := strings.Fields(line)
			if !strings.HasPrefix(field[10], "/dev/") {
				continue
			}
			diskUsageList = append(diskUsageList, DiskUsageInfo{
				Size:       field[0],
				Used:       field[1],
				Avail:      field[2],
				UsedPcent:  field[3],
				MountedOn:  field[4],
				Inodes:     field[5],
				IUsed:      field[6],
				IFree:      field[7],
				IUsedPcent: field[8],
				Type:       field[9],
				Filesystem: strings.Join(field[10:], " "),
			})
		}
	}

	result.Ok(c, diskUsageList)
}

// DiskIOStat info
// iostat -xdk
func (h *Handler) DiskIOStat(c *gin.Context) {
	var diskIOList [][]string
	var header []string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("iostat", "-xdk")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		lines := strings.Split(ClearNewline(stdout.String()), "\n")
		header = strings.Fields(lines[2])
		for _, line := range lines[3:] {
			diskIOList = append(diskIOList, strings.Fields(line))
		}
	}

	result.Ok(c, struct {
		Header []string   `json:"header"`
		List   [][]string `json:"list"`
	}{
		Header: header,
		List:   diskIOList,
	})
}

func ClearNewline(str string) string {
	return strings.TrimRight(strings.Replace(str, "\r\n", "\n", -1), "\n")
}

// @Tags Dashboard
// @Summary Load os info
// @Description 获取服务器基础数据
// @Accept json
// @Success 200 {object} dto.OsInfo
// @Security ApiKeyAuth
// @Router /dashboard/base/os [get]
func (h *Handler) LoadDashboardOsInfo(c *gin.Context) {
	data, err := h.Service.Dashboard.LoadOsInfo()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, data)
}

// @Tags Dashboard
// @Summary Load dashboard base info
// @Description 获取首页基础数据
// @Accept json
// @Param ioOption path string true "request"
// @Param netOption path string true "request"
// @Success 200 {object} dto.DashboardBase
// @Security ApiKeyAuth
// @Router /dashboard/base/:ioOption/:netOption [get]
func (h *Handler) LoadDashboardBaseInfo(c *gin.Context) {
	ioOption, ok := c.Params.Get("ioOption")
	if !ok {
		result.Error(c, "ioOption")
		return
	}
	netOption, ok := c.Params.Get("netOption")
	if !ok {
		result.Error(c, "netOption")
		return
	}
	data, err := h.Service.Dashboard.LoadBaseInfo(ioOption, netOption)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, data)
}

// @Tags Dashboard
// @Summary Load dashboard current info
// @Description 获取首页实时数据
// @Accept json
// @Param ioOption path string true "request"
// @Param netOption path string true "request"
// @Success 200 {object} dto.DashboardCurrent
// @Security ApiKeyAuth
// @Router /dashboard/current/:ioOption/:netOption [get]
func (h *Handler) LoadDashboardCurrentInfo(c *gin.Context) {
	ioOption, ok := c.Params.Get("ioOption")
	if !ok {
		result.Error(c, "ioOption")
		return
	}
	netOption, ok := c.Params.Get("netOption")
	if !ok {
		result.Error(c, "netOption")
		return
	}

	data := h.Service.Dashboard.LoadCurrentInfo(ioOption, netOption)
	result.Ok(c, data)
}
