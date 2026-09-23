package ai

import (
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/testutil"
)

// TestEffCfg_DBOverridesYaml 覆盖优先级：DB 非空字段覆盖 yaml，空字段回落。
func TestEffCfg_DBOverridesYaml(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewService(db)

	if err := config.StoreForTest(map[string]any{
		"ai.apiKey": "yaml-key", "ai.model": "yaml-model",
		"ai.baseURL": "https://yaml", "ai.maxTokens": 1000, "ai.timeout": 60,
		"ai.models": []map[string]any{{"name": "yamlm", "model": "ym"}},
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = config.StoreForTest(map[string]any{"ai.models": []map[string]any{}})
	})

	// 无覆盖行：全部来自 yaml
	eff := s.effCfg()
	if eff.Model != "yaml-model" || eff.BaseURL != "https://yaml" || len(eff.Models) != 1 {
		t.Fatalf("基线: %+v", eff)
	}

	// 写覆盖行：部分字段覆盖，models 替换，timeout 留 0 回落
	if err := s.SaveAIConfig(SaveAIConfigRequest{
		Model:     "db-model",
		APIKey:    "db-key",
		Models:    []ModelEntry{{Name: "dbm", Model: "dm", BaseURL: "https://db"}},
		MaxTokens: 4096,
	}); err != nil {
		t.Fatal(err)
	}
	eff = s.effCfg()
	if eff.Model != "db-model" || eff.APIKey != "db-key" || eff.MaxTokens != 4096 {
		t.Errorf("覆盖未生效: %+v", eff)
	}
	if eff.BaseURL != "https://yaml" || eff.Timeout != 60 {
		t.Errorf("空覆盖应回落 yaml: %+v", eff)
	}
	if len(eff.Models) != 1 || eff.Models[0].Name != "dbm" {
		t.Errorf("models 应被 DB 替换: %+v", eff.Models)
	}

	// Enabled 跟随覆盖（yaml key 清掉后靠 DB key）
	_ = config.StoreForTest(map[string]any{"ai.apiKey": ""})
	if !s.Enabled() {
		t.Error("DB 覆盖的 key 应使能")
	}
}

// TestConfigInfo_MaskedKey 接口视图：key 永不明文，掩码形态固定。
func TestConfigInfo_MaskedKey(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewService(db)

	if err := s.SaveAIConfig(SaveAIConfigRequest{APIKey: "glmc-1234567890abcdef"}); err != nil {
		t.Fatal(err)
	}
	info := s.ConfigInfo()
	if !info.APIKeySet {
		t.Fatal("apiKeySet 应为 true")
	}
	if strings.Contains(info.APIKeyMasked, "1234567890abcdef") {
		t.Error("掩码不得包含完整 key")
	}
	if !strings.HasPrefix(info.APIKeyMasked, "glmc") {
		t.Errorf("掩码应保留前 4 位: %q", info.APIKeyMasked)
	}
	if !info.FromDB["apiKey"] {
		t.Error("来源标注错误")
	}
}

// TestSaveAIConfig_KeySemantics key 三态语义：空=保持、清除、新值。
func TestSaveAIConfig_KeySemantics(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewService(db)

	// 设初值
	if err := s.SaveAIConfig(SaveAIConfigRequest{APIKey: "first-key"}); err != nil {
		t.Fatal(err)
	}
	// 空且未清除 → 保持
	if err := s.SaveAIConfig(SaveAIConfigRequest{Model: "m"}); err != nil {
		t.Fatal(err)
	}
	if got := s.loadAIConfigRow().APIKey; got != "first-key" {
		t.Errorf("空 key 应保持原值, got %q", got)
	}
	// 新值 → 覆盖
	if err := s.SaveAIConfig(SaveAIConfigRequest{APIKey: "second-key"}); err != nil {
		t.Fatal(err)
	}
	if got := s.loadAIConfigRow().APIKey; got != "second-key" {
		t.Errorf("got %q", got)
	}
	// 显式清除
	if err := s.SaveAIConfig(SaveAIConfigRequest{ClearAPIKey: true}); err != nil {
		t.Fatal(err)
	}
	if got := s.loadAIConfigRow().APIKey; got != "" {
		t.Errorf("清除后应为空, got %q", got)
	}
}

// TestSaveAIConfig_ModelValidation 模型条目校验：缺字段/重名拒绝。
func TestSaveAIConfig_ModelValidation(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewService(db)
	if err := s.SaveAIConfig(SaveAIConfigRequest{
		Models: []ModelEntry{{Name: "a"}}, // 缺 model
	}); err == nil {
		t.Error("缺 model 应拒绝")
	}
	if err := s.SaveAIConfig(SaveAIConfigRequest{
		Models: []ModelEntry{
			{Name: "a", Model: "m1"}, {Name: "a", Model: "m2"},
		},
	}); err == nil {
		t.Error("重名应拒绝")
	}
	if err := s.SaveAIConfig(SaveAIConfigRequest{
		DailyQuota: -1,
	}); err == nil {
		t.Error("负配额应拒绝")
	}
}
