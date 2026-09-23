package ai

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gitee.com/jieepre/keepblog/config"
)

// errInvalidConfig 配置值非法（负数配额/模型条目缺字段或重名）。
var errInvalidConfig = errors.New("配置值非法：请检查数值范围与模型条目（name/model 必填且不重复）")

// AIConfigRow 后台可运营的 AI 配置覆盖层（单行表，id=1）。
// 语义：字段留空/为 0 表示"回落 config.yaml 基线"，非空覆盖基线。
// apiKey 属敏感字段：接口只回掩码，永不明文出站。
type AIConfigRow struct {
	Id         int64 `gorm:"primaryKey"`
	BaseURL    string
	APIKey     string
	Model      string
	MaxTokens  int
	DailyQuota int
	StyleHint  string
	Timeout    int
	// ModelsJSON 多模型列表的 JSON 序列化；空串表示回落 yaml 的 models
	ModelsJSON string
	UpdatedAt  int64
}

const configTableDDL = `CREATE TABLE IF NOT EXISTS ai_config (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	base_url TEXT DEFAULT '',
	api_key TEXT DEFAULT '',
	model TEXT DEFAULT '',
	max_tokens INTEGER DEFAULT 0,
	daily_quota INTEGER DEFAULT 0,
	style_hint TEXT DEFAULT '',
	timeout INTEGER DEFAULT 0,
	models_json TEXT DEFAULT '',
	updated_at INTEGER
)`

// ModelEntry 多模型条目（服务内部表示，避免依赖 config 未导出类型）。
type ModelEntry struct {
	Name      string `json:"name"`
	Model     string `json:"model"`
	BaseURL   string `json:"baseURL"`
	APIKey    string `json:"apiKey"`
	MaxTokens int    `json:"maxTokens"`
}

// resolvedAI 合并后的生效配置。
type resolvedAI struct {
	BaseURL    string
	APIKey     string
	Model      string
	MaxTokens  int
	DailyQuota int
	StyleHint  string
	Timeout    int
	Models     []ModelEntry
}

// loadAIConfigRow 读覆盖行，无表/无行返回零值。
func (s *Service) loadAIConfigRow() AIConfigRow {
	var row AIConfigRow
	if s.db == nil {
		return row
	}
	if err := s.ensureTable(); err != nil {
		return row
	}
	// 无行时 Take 报 ErrRecordNotFound，零值即"全部回落"
	_ = s.db.Table("ai_config").Where("id = ?", 1).Take(&row).Error
	return row
}

// effCfg 生效配置 = yaml 基线 + DB 覆盖。每次现读，保存即时生效。
func (s *Service) effCfg() resolvedAI {
	var res resolvedAI
	if base := config.Get().Ai; base != nil {
		res.BaseURL, res.APIKey, res.Model = base.BaseURL, base.APIKey, base.Model
		res.MaxTokens, res.DailyQuota, res.StyleHint, res.Timeout =
			base.MaxTokens, base.DailyQuota, base.StyleHint, base.Timeout
		for _, m := range base.Models {
			res.Models = append(res.Models, ModelEntry{
				Name: m.Name, Model: m.Model, BaseURL: m.BaseURL,
				APIKey: m.APIKey, MaxTokens: m.MaxTokens,
			})
		}
	}
	row := s.loadAIConfigRow()
	if row.BaseURL != "" {
		res.BaseURL = row.BaseURL
	}
	if row.APIKey != "" {
		res.APIKey = row.APIKey
	}
	if row.Model != "" {
		res.Model = row.Model
	}
	if row.MaxTokens > 0 {
		res.MaxTokens = row.MaxTokens
	}
	if row.DailyQuota > 0 {
		res.DailyQuota = row.DailyQuota
	}
	if row.StyleHint != "" {
		res.StyleHint = row.StyleHint
	}
	if row.Timeout > 0 {
		res.Timeout = row.Timeout
	}
	if row.ModelsJSON != "" {
		var models []ModelEntry
		if err := json.Unmarshal([]byte(row.ModelsJSON), &models); err == nil && len(models) > 0 {
			res.Models = models
		}
	}
	return res
}

// AIConfigInfo GET 返回体：apiKey 只出掩码与已配置标记。
type AIConfigInfo struct {
	BaseURL      string       `json:"baseURL"`
	APIKeySet    bool         `json:"apiKeySet"`
	APIKeyMasked string       `json:"apiKeyMasked"`
	Model        string       `json:"model"`
	MaxTokens    int          `json:"maxTokens"`
	DailyQuota   int          `json:"dailyQuota"`
	StyleHint    string       `json:"styleHint"`
	Timeout      int          `json:"timeout"`
	Models       []ModelEntry `json:"models"`
	// FromDB 各字段是否来自后台覆盖（true）还是 yaml 基线（false），
	// 前端据此标注来源
	FromDB map[string]bool `json:"fromDB"`
}

// ConfigInfo 组装当前生效配置视图（含来源标注）。
func (s *Service) ConfigInfo() AIConfigInfo {
	row := s.loadAIConfigRow()
	eff := s.effCfg()

	info := AIConfigInfo{
		BaseURL:      eff.BaseURL,
		APIKeySet:    eff.APIKey != "",
		APIKeyMasked: maskKey(eff.APIKey),
		Model:        eff.Model,
		MaxTokens:    eff.MaxTokens,
		DailyQuota:   eff.DailyQuota,
		StyleHint:    eff.StyleHint,
		Timeout:      eff.Timeout,
		Models:       eff.Models,
		FromDB: map[string]bool{
			"baseURL":    row.BaseURL != "",
			"model":      row.Model != "",
			"maxTokens":  row.MaxTokens > 0,
			"dailyQuota": row.DailyQuota > 0,
			"styleHint":  row.StyleHint != "",
			"timeout":    row.Timeout > 0,
			"models":     row.ModelsJSON != "",
			"apiKey":     row.APIKey != "",
		},
	}
	return info
}

// maskKey 前 4 后 4，中间打码；短 key 全打码。
func maskKey(key string) string {
	if key == "" {
		return ""
	}
	r := []rune(strings.TrimSpace(key))
	if len(r) <= 8 {
		return strings.Repeat("*", len(r))
	}
	return string(r[:4]) + strings.Repeat("*", 6) + string(r[len(r)-4:])
}

// SaveAIConfigRequest PUT 请求体。语义：
//   - 普通字段直接写覆盖层——空值即清除该字段覆盖（回落 yaml）
//   - apiKey 空串 = 保持已存覆盖不变；ClearAPIKey=true 显式清除；
//     非空 = 写入新值
type SaveAIConfigRequest struct {
	BaseURL     string       `json:"baseURL"`
	APIKey      string       `json:"apiKey"`
	ClearAPIKey bool         `json:"clearApiKey"`
	Model       string       `json:"model"`
	MaxTokens   int          `json:"maxTokens"`
	DailyQuota  int          `json:"dailyQuota"`
	StyleHint   string       `json:"styleHint"`
	Timeout     int          `json:"timeout"`
	Models      []ModelEntry `json:"models"`
}

// SaveAIConfig 保存覆盖层。models 内部条目要求 name/model 非空且 name 不重复。
func (s *Service) SaveAIConfig(req SaveAIConfigRequest) error {
	if err := s.ensureTable(); err != nil {
		return err
	}
	if req.Timeout < 0 || req.MaxTokens < 0 || req.DailyQuota < 0 {
		return errInvalidConfig
	}
	seen := make(map[string]bool)
	models := make([]ModelEntry, 0, len(req.Models))
	for _, m := range req.Models {
		if strings.TrimSpace(m.Name) == "" || strings.TrimSpace(m.Model) == "" {
			return errInvalidConfig
		}
		if seen[m.Name] {
			return errInvalidConfig
		}
		seen[m.Name] = true
		models = append(models, m)
	}

	row := s.loadAIConfigRow()
	row.Id = 1
	row.BaseURL = strings.TrimSpace(req.BaseURL)
	row.Model = strings.TrimSpace(req.Model)
	row.MaxTokens = req.MaxTokens
	row.DailyQuota = req.DailyQuota
	row.StyleHint = strings.TrimSpace(req.StyleHint)
	row.Timeout = req.Timeout
	if len(models) > 0 {
		data, err := json.Marshal(models)
		if err != nil {
			return err
		}
		row.ModelsJSON = string(data)
	} else {
		row.ModelsJSON = ""
	}
	switch {
	case req.ClearAPIKey:
		row.APIKey = ""
	case strings.TrimSpace(req.APIKey) != "":
		row.APIKey = strings.TrimSpace(req.APIKey)
		// 空 + 未勾选清除 = 保持原值（row.APIKey 未动）
	}
	row.UpdatedAt = time.Now().Unix()
	return s.db.Table("ai_config").Save(&row).Error
}
