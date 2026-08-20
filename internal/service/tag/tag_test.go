package tag

import (
	"testing"

	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/testutil"
)

func TestGetTagList(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewTagService()

	list, err := s.GetTagList()
	if err != nil {
		t.Fatalf("GetTagList 报错: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("len = %d, want 3", len(list))
	}
}

func TestSave_IdempotentByTagName(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewTagService()

	// 已存在的标签名返回既有 id，不重复创建
	id, err := s.Save(model.Tag{TagName: "go"})
	if err != nil {
		t.Fatalf("Save 报错: %v", err)
	}
	if id != 1 {
		t.Errorf("Save(go) 返回 id = %d, want 1", id)
	}

	newId, err := s.Save(model.Tag{TagName: "全新标签"})
	if err != nil {
		t.Fatalf("Save(新标签) 报错: %v", err)
	}
	if newId == 0 {
		t.Error("新标签应返回非零 id")
	}
}

func TestGetTag(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewTagService()

	tag, err := s.GetTag(2)
	if err != nil {
		t.Fatalf("GetTag 报错: %v", err)
	}
	if tag.TagName != "mysql" {
		t.Errorf("TagName = %q, want mysql", tag.TagName)
	}
	if _, err := s.GetTag(999); err == nil {
		t.Error("不存在的标签应报错")
	}
}

func TestGetTags_OnlyPublishedPosts(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewTagService()

	tags, err := s.GetTags()
	if err != nil {
		t.Fatalf("GetTags 报错: %v", err)
	}
	// 已删除文章 P91 上的 go 标签不应使结果去重异常；共 3 个标签
	if len(tags) != 3 {
		t.Fatalf("len = %d, want 3: %+v", len(tags), tags)
	}
	names := map[string]bool{}
	for _, tag := range tags {
		names[tag.TagName] = true
		if tag.TagStyle == "" {
			t.Errorf("标签 %q 缺少云样式", tag.TagName)
		}
	}
	for _, want := range []string{"go", "mysql", "sqlite"} {
		if !names[want] {
			t.Errorf("缺少标签 %q", want)
		}
	}
}

func TestUpdate(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewTagService()

	if err := s.Update(model.Tag{TagId: 3, TagName: "sqlite3"}); err != nil {
		t.Fatalf("Update 报错: %v", err)
	}
	tag, _ := s.GetTag(3)
	if tag.TagName != "sqlite3" {
		t.Errorf("TagName = %q, want sqlite3", tag.TagName)
	}
}

func TestDelete_RemovesAssociations(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewTagService()

	// 自建标签与关联，避免破坏共享 fixture
	id, err := s.Save(model.Tag{TagName: "待删除标签"})
	if err != nil {
		t.Fatalf("Save 报错: %v", err)
	}
	db := testutil.NewTestDB(t)
	if err := db.Create(&model.PostTag{PostId: 1, TagId: id}).Error; err != nil {
		t.Fatalf("建立关联失败: %v", err)
	}

	if err := s.Delete(id); err != nil {
		t.Fatalf("Delete 报错: %v", err)
	}
	if _, err := s.GetTag(id); err == nil {
		t.Error("标签应已删除")
	}
	var count int64
	db.Model(&model.PostTag{}).Where("tag_id = ?", id).Count(&count)
	if count != 0 {
		t.Errorf("post_tag 关联残留 %d 条", count)
	}
}
