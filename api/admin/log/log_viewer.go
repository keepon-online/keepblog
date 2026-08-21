package log

import (
	"gitee.com/jieepre/keepblog/internal/logger"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

// LogLevelRequest 日志级别请求
type LogLevelRequest struct {
	Level string `json:"level" binding:"required,oneof=debug info warn error"`
}

// LogQueryRequest 日志查询请求
type LogQueryRequest struct {
	Filename string `form:"filename"`
	Lines    int    `form:"lines" default:"100"`
	Search   string `form:"search"`
}

// GetLogLevel 获取当前日志级别
// @Summary 获取日志级别
// @Tags 日志管理
// @Success 200 {object} result.Response
// @Router /api/log/level [get]
func (h *Handler) GetLogLevel(c *gin.Context) {
	level := logger.GetLevel()
	config := logger.GetConfig()

	result.Ok(c, gin.H{
		"level":  level,
		"format": config.Format,
		"output": config.Output,
		"path":   config.Path,
	})
}

// SetLogLevel 动态设置日志级别
// @Summary 设置日志级别
// @Tags 日志管理
// @Param data body LogLevelRequest true "日志级别"
// @Success 200 {object} result.Response
// @Router /api/log/level [put]
func (h *Handler) SetLogLevel(c *gin.Context) {
	var req LogLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, "无效的日志级别，可选值: debug, info, warn, error")
		return
	}

	if err := logger.SetLevel(req.Level); err != nil {
		result.Error(c, err.Error())
		return
	}

	result.Ok(c, gin.H{
		"message": "日志级别已更新",
		"level":   req.Level,
	})
}

// GetLogStats 获取日志统计信息
// @Summary 获取日志统计
// @Tags 日志管理
// @Success 200 {object} result.Response
// @Router /api/log/stats [get]
func (h *Handler) GetLogStats(c *gin.Context) {
	stats, err := logger.GetLogStats()
	if err != nil {
		result.Error(c, err.Error())
		return
	}

	result.Ok(c, stats)
}

// GetLogList 获取日志文件列表
// @Summary 获取日志列表
// @Tags 日志管理
// @Success 200 {object} result.Response
// @Router /api/log/list [get]
func (h *Handler) GetLogList(c *gin.Context) {
	stats, err := logger.GetLogStats()
	if err != nil {
		result.Error(c, err.Error())
		return
	}

	result.Ok(c, gin.H{
		"files": stats.LogFiles,
		"total": stats.TotalFiles,
	})
}

// ReadLog 读取日志内容
// @Summary 读取日志内容
// @Tags 日志管理
// @Param filename query string true "日志文件名"
// @Param lines query int false "返回行数" default(100)
// @Param search query string false "搜索关键词"
// @Success 200 {object} result.Response
// @Router /api/log/read [get]
func (h *Handler) ReadLog(c *gin.Context) {
	var req LogQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Error(c, "参数错误")
		return
	}

	if req.Filename == "" {
		result.Error(c, "请指定日志文件名")
		return
	}

	if req.Lines <= 0 {
		req.Lines = 100
	}

	lines, err := logger.ReadLogFile(req.Filename, req.Lines, req.Search)
	if err != nil {
		result.Error(c, err.Error())
		return
	}

	result.Ok(c, gin.H{
		"filename": req.Filename,
		"lines":    lines,
		"count":    len(lines),
	})
}
