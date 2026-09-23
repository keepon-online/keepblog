// Package ai 后台 AI 辅助写作服务：prompt 组装、每日配额与用量审计。
// 内容安全原则：用量表只记任务类型与 token 数，不落任何文章内容。
package ai

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/pkg/ai"
	"gorm.io/gorm"
)

// 支持的任务与润色模式。新增任务时同步 BuildMessages 与前端入口。
const (
	TaskPolish    = "polish"
	TaskContinue  = "continue"
	TaskTitle     = "title"
	TaskSummary   = "summary"
	TaskRefine    = "refine"    // 对上一版结果按追加指令再生成
	TaskTags      = "tags"      // 全文标签建议
	TaskProofread = "proofread" // 全文校对，输出问题清单
)

// EditRequest AI 编辑请求。上下文片段由前端采集：
// Selection 为选中文本；Before/After 为选区前后文窗口（润色），
// 续写任务里 Before/After 语义为光标前后窗口（有 Before 优先于 Digest 尾部截断）；
// Digest 为长文压缩摘要（续写/标题/摘要/标签/校对任务用全文时前端先行截断）；
// Refine 任务：Previous 为上一版结果，Instruction 为追加修饰要求。
type EditRequest struct {
	Task string `json:"task" binding:"required"`
	Mode string `json:"mode"`
	// Model 多模型切换：ai.models 里的 name，空用默认模型
	Model       string `json:"model"`
	Selection   string `json:"selection"`
	Before      string `json:"before"`
	After       string `json:"after"`
	Title       string `json:"title"`
	Series      string `json:"series"`
	Digest      string `json:"digest"`
	Previous    string `json:"previous"`
	Instruction string `json:"instruction"`
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

const templateTableDDL = `CREATE TABLE IF NOT EXISTS ai_template (
	tpl_key TEXT PRIMARY KEY,
	content TEXT DEFAULT '',
	updated_at INTEGER
)`

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

// EffectiveStyleHint 生效配置的文风设定（prompt 组装用）。
func (s *Service) EffectiveStyleHint() string {
	return s.effCfg().StyleHint
}

// Enabled 返回 AI 功能是否可用（生效配置里 apiKey 非空，yaml 或后台覆盖均可）。
func (s *Service) Enabled() bool {
	return strings.TrimSpace(s.effCfg().APIKey) != ""
}

// ModelOption 多模型切换的可选项。
type ModelOption struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

// StatusInfo /api/v1/ai/status 返回体。
type StatusInfo struct {
	Enabled      bool          `json:"enabled"`
	Model        string        `json:"model"`
	TodayUsed    int           `json:"todayUsed"`
	DailyQuota   int           `json:"dailyQuota"`
	DefaultModel string        `json:"defaultModel"`
	Models       []ModelOption `json:"models"`
}

// Status 探测接口数据：未配置时只返回 enabled=false。
func (s *Service) Status() StatusInfo {
	if !s.Enabled() {
		return StatusInfo{Enabled: false}
	}
	eff := s.effCfg()
	return StatusInfo{
		Enabled:      true,
		Model:        eff.Model,
		TodayUsed:    s.todayCount(),
		DailyQuota:   eff.DailyQuota,
		DefaultModel: eff.Model,
		Models:       s.modelOptions(),
	}
}

// ErrQuotaExceeded 每日配额用尽。
var ErrQuotaExceeded = errors.New("今日 AI 调用配额已用完")

// AllowCall 配额检查（调用前）。
func (s *Service) AllowCall() error {
	if !s.Enabled() {
		return errors.New("AI 未配置（config.ai.apiKey 为空）")
	}
	eff := s.effCfg()
	if eff.DailyQuota <= 0 {
		return nil // 0 = 不限制
	}
	if s.todayCount() >= eff.DailyQuota {
		return ErrQuotaExceeded
	}
	return nil
}

// Client 按生效配置（yaml+后台覆盖）构造上游客户端。
func (s *Service) Client() *ai.Client {
	eff := s.effCfg()
	return ai.New(eff.BaseURL, eff.APIKey, eff.Model, eff.Timeout)
}

// ErrModelNotFound 请求的模型名未在 ai.models 配置中。
var ErrModelNotFound = errors.New("未配置的模型名：")

// resolveModel 按 name 解析实际上游参数：空 name 用顶层默认；配置了
// models 且命中则条目字段逐项覆盖（缺省继承顶层）。未命中返回错误而
// 非静默回退默认——用户以为在用 A 模型实际用默认是不可接受的。
func (s *Service) resolveModel(name string) (baseURL, apiKey, model string, maxTokens int, err error) {
	eff := s.effCfg()
	baseURL, apiKey, model, maxTokens = eff.BaseURL, eff.APIKey, eff.Model, eff.MaxTokens
	if name == "" {
		return baseURL, apiKey, model, maxTokens, nil
	}
	if len(eff.Models) > 0 {
		for _, m := range eff.Models {
			if m.Name == name {
				if m.BaseURL != "" {
					baseURL = m.BaseURL
				}
				if m.APIKey != "" {
					apiKey = m.APIKey
				}
				if m.Model != "" {
					model = m.Model
				}
				if m.MaxTokens > 0 {
					maxTokens = m.MaxTokens
				}
				return baseURL, apiKey, model, maxTokens, nil
			}
		}
	}
	// name 恰好等于默认模型名也放行（前端"默认"选项不带 name 的兜底）
	if name == eff.Model {
		return baseURL, apiKey, model, maxTokens, nil
	}
	return "", "", "", 0, fmt.Errorf("%w%s", ErrModelNotFound, name)
}

// ClientFor 按请求指定的模型名构造客户端。
func (s *Service) ClientFor(name string) (*ai.Client, error) {
	baseURL, apiKey, model, _, err := s.resolveModel(name)
	if err != nil {
		return nil, err
	}
	return ai.New(baseURL, apiKey, model, timeoutOf(s.effCfg().Timeout)), nil
}

// timeoutOf 超时兜底（<=0 用默认 120s）。
func timeoutOf(sec int) int {
	if sec <= 0 {
		return 120
	}
	return sec
}

// MaxTokensFor 按请求指定的模型名取单次生成上限。
func (s *Service) MaxTokensFor(name string) int {
	if _, _, _, mt, err := s.resolveModel(name); err == nil && mt > 0 {
		return mt
	}
	return s.MaxTokens()
}

// ActualModel 请求名对应的真实模型标识（用量记录显示用）。
func (s *Service) ActualModel(name string) string {
	_, _, model, _, err := s.resolveModel(name)
	if err != nil {
		return s.effCfg().Model
	}
	return model
}

// modelOptions 状态接口的可选模型列表（默认恒在首位）。
// 不导出 config 的 ai 类型，方法内直接读快照。
func (s *Service) modelOptions() []ModelOption {
	eff := s.effCfg()
	opts := []ModelOption{{Name: "default", Model: eff.Model}}
	for _, m := range eff.Models {
		opts = append(opts, ModelOption{Name: m.Name, Model: m.Model})
	}
	return opts
}

// TemplateRow 模板存储行。key 列用 tpl_key 避开 SQL 关键字。
type TemplateRow struct {
	Key       string `gorm:"column:tpl_key;primaryKey"`
	Content   string
	UpdatedAt int64
}

// TemplateOverride 取任务的模板覆盖，未配置返回空串。
func (s *Service) TemplateOverride(key string) string {
	if s.db == nil {
		return ""
	}
	if err := s.ensureTable(); err != nil {
		return ""
	}
	var row TemplateRow
	if err := s.db.Table("ai_template").
		Where("tpl_key = ?", key).
		Take(&row).Error; err != nil {
		return ""
	}
	return row.Content
}

// SaveTemplate 保存/清除模板覆盖（content 为空即删除恢复默认）。
// key 白名单校验防写入任意键。
func (s *Service) SaveTemplate(key, content string) error {
	if !IsTemplateKey(key) {
		return errors.New("不支持的任务键: " + key)
	}
	if len([]rune(content)) > templateLimitRune() {
		return errors.New("模板过长")
	}
	if err := s.ensureTable(); err != nil {
		return err
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return s.db.Table("ai_template").Where("tpl_key = ?", key).
			Delete(&TemplateRow{}).Error
	}
	return s.db.Table("ai_template").Save(&TemplateRow{
		Key: key, Content: content, UpdatedAt: time.Now().Unix(),
	}).Error
}

func templateLimitRune() int { return 500 }

// TemplateItems 模板管理列表：内置默认 + 已保存覆盖。
func (s *Service) TemplateItems() []TemplateItem {
	saved := make(map[string]string)
	if s.db != nil {
		if err := s.ensureTable(); err == nil {
			var rows []TemplateRow
			if err := s.db.Table("ai_template").Find(&rows).Error; err == nil {
				for _, r := range rows {
					saved[r.Key] = r.Content
				}
			}
		}
	}
	items := DefaultTemplateItems()
	for i := range items {
		items[i].Content = saved[items[i].Key]
	}
	return items
}

// MaxTokens 单次生成上限（0 时用安全默认值）。
func (s *Service) MaxTokens() int {
	if cfg := config.Get().Ai; cfg != nil && cfg.MaxTokens > 0 {
		return cfg.MaxTokens
	}
	return 2048
}

// RecordUsage 异步记录一次调用。model 为实际使用的模型标识。
// 失败只影响统计不影响主流程。
func (s *Service) RecordUsage(task, model string, usage ai.Usage) {
	if s.db == nil {
		return
	}
	row := UsageRow{
		Task: task, Model: model,
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
		if err := s.db.Exec(usageTableDDL).Error; err != nil {
			s.ensureErr = err
			return
		}
		if err := s.db.Exec(templateTableDDL).Error; err != nil {
			s.ensureErr = err
			return
		}
		s.ensureErr = s.db.Exec(configTableDDL).Error
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
