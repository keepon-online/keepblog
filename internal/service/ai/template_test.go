package ai

import (
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/testutil"
)

// TestBuildMessages_TemplateOverride 自定义模板替换指令位，
// 上下文组装与防注入 fence 保持不变。
func TestBuildMessages_TemplateOverride(t *testing.T) {
	msgs, err := BuildMessages(EditRequest{
		Task: TaskPolish, Selection: "原文", Before: "前", After: "后",
	}, "按我的个人风格改写，多用短句", "")
	if err != nil {
		t.Fatal(err)
	}
	user := msgs[1].Content
	if !strings.HasPrefix(user, "按我的个人风格改写，多用短句\n") {
		t.Errorf("指令位未使用模板: %s", user[:min(50, len(user))])
	}
	for _, need := range []string{"【选中文本】\n原文", "前", "后", fence} {
		if !strings.Contains(user, need) {
			t.Errorf("模板模式下丢失上下文 %q", need)
		}
	}

	// 超长模板截断
	long := strings.Repeat("字", 600)
	msgs, _ = BuildMessages(EditRequest{Task: TaskSummary, Digest: "x"}, long, "")
	if !strings.Contains(msgs[1].Content, "已截断") {
		t.Error("超长模板应截断")
	}
}

func TestTemplateKey(t *testing.T) {
	cases := map[[2]string]string{
		{"polish", ""}:         "polish:polish",
		{"polish", "expand"}:   "polish:expand",
		{"continue", ""}:       "continue",
		{"summary", ""}:        "summary",
		{"summary", "ignored"}: "summary",
	}
	for in, want := range cases {
		if got := TemplateKey(in[0], in[1]); got != want {
			t.Errorf("TemplateKey(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestTemplateStorage_SaveListDelete(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewService(db)

	if err := s.SaveTemplate("polish:polish", "自定义润色指令"); err != nil {
		t.Fatal(err)
	}
	if got := s.TemplateOverride("polish:polish"); got != "自定义润色指令" {
		t.Errorf("override = %q", got)
	}

	items := s.TemplateItems()
	found := false
	for _, it := range items {
		if it.Key == "polish:polish" {
			found = true
			if it.Content != "自定义润色指令" || it.DefaultContent == "" {
				t.Errorf("条目 = %+v", it)
			}
		}
	}
	if !found {
		t.Error("列表缺少已保存模板")
	}

	// 白名单拒绝
	if err := s.SaveTemplate("hack:evil", "x"); err == nil {
		t.Error("非白名单键应拒绝")
	}

	// 空内容删除恢复默认
	if err := s.SaveTemplate("polish:polish", ""); err != nil {
		t.Fatal(err)
	}
	if got := s.TemplateOverride("polish:polish"); got != "" {
		t.Errorf("删除后 override = %q, want 空", got)
	}
}

// TestResolveModel 多模型解析：默认/命中/继承/未知。
func TestResolveModel(t *testing.T) {
	db := testutil.NewTestDB(t)
	_ = db
	s := NewService(db)

	// 未配置 models：一切走默认
	if err := config.StoreForTest(map[string]any{
		"ai.apiKey": "k", "ai.model": "default-m", "ai.baseURL": "https://d",
	}); err != nil {
		t.Fatal(err)
	}
	_, _, model, _, err := s.resolveModel("")
	if err != nil || model != "default-m" {
		t.Fatalf("默认解析: %v %s", err, model)
	}
	if _, _, _, _, err := s.resolveModel("unknown"); err == nil {
		t.Error("未配置 models 时未知名应报错")
	}

	// 配置 models：命中覆盖、缺省继承、未知报错
	if err := config.StoreForTest(map[string]any{
		"ai.apiKey": "k", "ai.model": "default-m", "ai.baseURL": "https://d", "ai.maxTokens": 1000,
		"ai.models": []map[string]any{
			{"name": "fast", "model": "mini-m"},
			{"name": "ds", "model": "ds-chat", "baseURL": "https://api.ds", "apiKey": "ds-key", "maxTokens": 4096},
		},
	}); err != nil {
		t.Fatal(err)
	}
	baseURL, key, model, mt, err := s.resolveModel("fast")
	if err != nil || model != "mini-m" || baseURL != "https://d" || key != "k" || mt != 1000 {
		t.Errorf("fast: %s %s %s %d %v（应继承顶层）", baseURL, key, model, mt, err)
	}
	baseURL, key, model, mt, err = s.resolveModel("ds")
	if err != nil || model != "ds-chat" || baseURL != "https://api.ds" || key != "ds-key" || mt != 4096 {
		t.Errorf("ds: %s %s %s %d %v（应覆盖）", baseURL, key, model, mt, err)
	}
	if _, _, _, _, err := s.resolveModel("nope"); err == nil {
		t.Error("未知名应报错（不静默回退默认）")
	}
	// name 恰为默认模型名放行
	if _, _, m, _, err := s.resolveModel("default-m"); err != nil || m != "default-m" {
		t.Errorf("默认名直传: %v", err)
	}

	// 清理测试注入
	_ = config.StoreForTest(map[string]any{"ai.models": []map[string]any{}})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
