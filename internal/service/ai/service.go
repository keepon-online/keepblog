// Package ai 后台 AI 辅助写作服务：prompt 组装、每日配额与用量审计。
// 内容安全原则：用量表只记任务类型与 token 数，不落任何文章内容。
package ai

import (
	"errors"
	"strings"
	"sync"
	"time"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/pkg/ai"
	"gorm.io/gorm"
)

// 支持的任务与润色模式。新增任务时同步 BuildMessages 与前端入口。
const (
	TaskPolish   = "polish"
	TaskContinue = "continue"
	TaskTitle    = "title"
	TaskSummary  = "summary"
)

// EditRequest AI 编辑请求。上下文片段由前端采集：
// Selection 为选中文本；Before/After 为选区前后文窗口（润色）；
// Digest 为长文压缩摘要（续写/标题/摘要任务用全文时前端先行截断）。
type EditRequest struct {
	Task      string `json:"task" binding:"required"`
	Mode      string `json:"mode"`
	Selection string `json:"selection"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Title     string `json:"title"`
	Series    string `json:"series"`
	Digest    string `json:"digest"`
}

// UsageRow 用量审计表。按 ensure-column 同款思路：建表语句幂等，
// 不进 goose 迁移（避免与 AutoMigrate 引导路径冲突）。
type UsageRow struct {
	Id               uint64 `gorm:"primaryKey;autoIncrement"`
	Task             string `gorm:"default:''"`
	Model            string `gorm:"default:''"`
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	CreatedAt        int64 `gorm:"autoCreateTime"`
}

const usageTableDDL = `CREATE TABLE IF NOT EXISTS ai_usage (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	task TEXT DEFAULT '',
	model TEXT DEFAULT '',
	prompt_tokens INTEGER DEFAULT 0,
	completion_tokens INTEGER DEFAULT 0,
	total_tokens INTEGER DEFAULT 0,
	created_at INTEGER
)`

// Service AI 服务。db 为 nil 时仅禁用状态查询可用（防御性，正常装配不会出现）。
type Service struct {
	db *gorm.DB
	// 建表一次的守卫；db 为 nil 时直接短路
	ensureOnce sync.Once
	ensureErr  error
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Enabled 返回 AI 功能是否可用（配置了 apiKey）。
func (s *Service) Enabled() bool {
	cfg := config.Get().Ai
	return cfg != nil && strings.TrimSpace(cfg.APIKey) != ""
}

// StatusInfo /api/v1/ai/status 返回体。
type StatusInfo struct {
	Enabled    bool   `json:"enabled"`
	Model      string `json:"model"`
	TodayUsed  int    `json:"todayUsed"`
	DailyQuota int    `json:"dailyQuota"`
}

// Status 探测接口数据：未配置时只返回 enabled=false。
func (s *Service) Status() StatusInfo {
	if !s.Enabled() {
		return StatusInfo{Enabled: false}
	}
	cfg := config.Get().Ai
	return StatusInfo{
		Enabled:    true,
		Model:      cfg.Model,
		TodayUsed:  s.todayCount(),
		DailyQuota: cfg.DailyQuota,
	}
}

// ErrQuotaExceeded 每日配额用尽。
var ErrQuotaExceeded = errors.New("今日 AI 调用配额已用完")

// AllowCall 配额检查（调用前）。
func (s *Service) AllowCall() error {
	if !s.Enabled() {
		return errors.New("AI 未配置（config.ai.apiKey 为空）")
	}
	cfg := config.Get().Ai
	if cfg.DailyQuota <= 0 {
		return nil // 0 = 不限制
	}
	if s.todayCount() >= cfg.DailyQuota {
		return ErrQuotaExceeded
	}
	return nil
}

// Client 按当前配置构造上游客户端（配置热重载后下次调用即生效）。
func (s *Service) Client() *ai.Client {
	cfg := config.Get().Ai
	return ai.New(cfg.BaseURL, cfg.APIKey, cfg.Model, cfg.Timeout)
}

// MaxTokens 单次生成上限（0 时用安全默认值）。
func (s *Service) MaxTokens() int {
	if cfg := config.Get().Ai; cfg != nil && cfg.MaxTokens > 0 {
		return cfg.MaxTokens
	}
	return 2048
}

// RecordUsage 异步记录一次调用。失败只影响统计不影响主流程。
func (s *Service) RecordUsage(task string, usage ai.Usage) {
	if s.db == nil {
		return
	}
	row := UsageRow{
		Task: task, Model: config.Get().Ai.Model,
		PromptTokens: usage.PromptTokens, CompletionTokens: usage.CompletionTokens,
		TotalTokens: usage.TotalTokens, CreatedAt: time.Now().Unix(),
	}
	go func() {
		if err := s.ensureTable(); err != nil {
			return
		}
		_ = s.db.Table("ai_usage").Create(&row).Error
	}()
}

func (s *Service) ensureTable() error {
	s.ensureOnce.Do(func() {
		if s.db == nil {
			return
		}
		s.ensureErr = s.db.Exec(usageTableDDL).Error
	})
	return s.ensureErr
}

// todayCount 今日已调用次数（表缺失时按 0 处理，宽松不阻断）。
func (s *Service) todayCount() int {
	if s.db == nil {
		return 0
	}
	if err := s.ensureTable(); err != nil {
		return 0
	}
	dayStart := time.Now().Truncate(24 * time.Hour).Unix()
	var n int64
	if err := s.db.Table("ai_usage").
		Where("created_at >= ?", dayStart).
		Count(&n).Error; err != nil {
		return 0
	}
	return int(n)
}
