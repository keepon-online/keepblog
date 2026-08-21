package about

import (
	"testing"

	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/testutil"
)

func TestInfo(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewAboutService(testutil.NewTestDB(t))

	info, err := s.Info()
	if err != nil {
		t.Fatalf("Info 报错: %v", err)
	}
	if info.Note != "# 关于博客" {
		t.Errorf("Note = %q, want %q", info.Note, "# 关于博客")
	}
}

func TestSave_Overwrite(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewAboutService(testutil.NewTestDB(t))

	// Save 对既有 Id 是整行覆盖
	if err := s.Save(model.About{Id: 1, Title: "新的标题", Note: "# 新内容"}); err != nil {
		t.Fatalf("Save 报错: %v", err)
	}
	info, _ := s.Info()
	if info.Note != "# 新内容" || info.Title != "新的标题" {
		t.Errorf("覆盖结果 = %+v", info)
	}
}
