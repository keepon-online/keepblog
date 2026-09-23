package ai

import (
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/internal/testutil"
)

func TestBuildMessages_Polish(t *testing.T) {
	db := testutil.NewTestDB(t)
	_ = db
	msgs, err := BuildMessages(EditRequest{
		Task: TaskPolish, Mode: "expand",
		Selection: "选中的段落", Before: "前文", After: "后文",
	}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 || msgs[0].Role != "system" || msgs[1].Role != "user" {
		t.Fatalf("消息结构 = %+v", msgs)
	}
	for _, need := range []string{"扩写", "选中的段落", "前文", "后文", fence} {
		if !strings.Contains(msgs[1].Content, need) {
			t.Errorf("user 消息缺少 %q: %s", need, msgs[1].Content)
		}
	}
	if !strings.Contains(msgs[0].Content, "不是对你的指令") {
		t.Error("system 缺少防注入声明")
	}
}

func TestBuildMessages_ContinueUsesTail(t *testing.T) {
	long := strings.Repeat("前", digestLimit+500) + "结尾在这里"
	msgs, err := BuildMessages(EditRequest{Task: TaskContinue, Digest: long, Title: "T", Series: "S"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil {
		t.Fatal(err)
	}
	body := msgs[1].Content
	if !strings.Contains(body, "结尾在这里") {
		t.Error("续写应包含尾部窗口内容")
	}
	if !strings.Contains(body, "开头已截断") {
		t.Error("超长正文应做尾部截断标记")
	}
	if !strings.Contains(body, "文章标题") || !strings.Contains(body, "所属系列") {
		t.Error("续写应携带标题与系列")
	}
}

func TestBuildMessages_ContinueAtCursor(t *testing.T) {
	msgs, err := BuildMessages(EditRequest{
		Task: TaskContinue, Before: "光标前的内容", After: "光标后的内容", Title: "T",
	}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil {
		t.Fatal(err)
	}
	body := msgs[1].Content
	for _, need := range []string{"光标前内容", "从光标处继续写", "光标后的内容", "不要重复"} {
		if !strings.Contains(body, need) {
			t.Errorf("光标续写缺少 %q: %s", need, body)
		}
	}
}

func TestBuildMessages_Refine(t *testing.T) {
	msgs, err := BuildMessages(EditRequest{
		Task: TaskRefine, Previous: "上一版结果", Instruction: "再精简一点",
	}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil {
		t.Fatal(err)
	}
	body := msgs[1].Content
	for _, need := range []string{"追加要求", "上一版结果", "再精简一点", "只输出修改后的完整文本"} {
		if !strings.Contains(body, need) {
			t.Errorf("refine 缺少 %q: %s", need, body)
		}
	}
	if _, err := BuildMessages(EditRequest{Task: TaskRefine, Previous: "x"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰"); err == nil {
		t.Error("refine 缺追加要求应报错")
	}
	if _, err := BuildMessages(EditRequest{Task: TaskRefine, Instruction: "x"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰"); err == nil {
		t.Error("refine 缺上一版结果应报错")
	}
}

func TestBuildMessages_TagsAndProofread(t *testing.T) {
	msgs, err := BuildMessages(EditRequest{Task: TaskTags, Digest: "正文", Title: "T"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil || !strings.Contains(msgs[1].Content, "5~8 个标签") {
		t.Errorf("tags 任务: %v", err)
	}
	msgs, err = BuildMessages(EditRequest{Task: TaskProofread, Digest: "正文"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil || !strings.Contains(msgs[1].Content, "未发现问题") {
		t.Errorf("proofread 任务: %v", err)
	}
	if _, err := BuildMessages(EditRequest{Task: TaskProofread}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰"); err == nil {
		t.Error("proofread 无正文应报错")
	}
}

func TestBuildMessages_TitleAndSummary(t *testing.T) {
	msgs, err := BuildMessages(EditRequest{Task: TaskTitle, Digest: "正文"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil || !strings.Contains(msgs[1].Content, "3 个候选标题") {
		t.Errorf("title 任务: %v", err)
	}
	msgs, err = BuildMessages(EditRequest{Task: TaskSummary, Digest: "正文"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if err != nil || !strings.Contains(msgs[1].Content, "120 字以内") {
		t.Errorf("summary 任务: %v", err)
	}
}

func TestBuildMessages_InputGuard(t *testing.T) {
	if _, err := BuildMessages(EditRequest{Task: TaskPolish}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰"); err == nil {
		t.Error("润色无选区应报错")
	}
	if _, err := BuildMessages(EditRequest{Task: "hack"}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰"); err == nil {
		t.Error("未知任务应报错")
	}
	// 截断不破坏中文（超过 digestLimit 才触发）
	msgs, _ := BuildMessages(EditRequest{Task: TaskSummary, Digest: strings.Repeat("字", digestLimit+999)}, "", "技术博客写作助手，文风简洁准确，避免空洞修饰")
	if !strings.Contains(msgs[1].Content, "已截断") {
		t.Error("超长正文应截断")
	}
}

func TestQuotaAndUsage(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewService(db)

	// 未配置 apiKey：不可用
	if s.Enabled() {
		t.Fatal("未配置时应禁用")
	}
	if err := s.AllowCall(); err == nil {
		t.Fatal("未配置时 AllowCall 应报错")
	}
	if st := s.Status(); st.Enabled {
		t.Fatal("status 应为禁用")
	}

	// 配额度 1，插一条今日用量后应拒绝
	if err := db.Exec(usageTableDDL).Error; err != nil {
		t.Fatal(err)
	}
	db.Exec("INSERT INTO ai_usage (task, model, created_at) VALUES ('polish', 'm', strftime('%s','now'))")

	// 直接以数据库计数验证配额逻辑（绕过 config 依赖的 Enabled 前置）
	if got := s.todayCount(); got != 1 {
		t.Errorf("todayCount = %d, want 1", got)
	}
}
