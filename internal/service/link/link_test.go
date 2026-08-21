package link

import (
	"testing"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/testutil"
)

func TestGetLinks_OnlyActive(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewLinkService(testutil.NewTestDB(t))

	links, err := s.GetLinks()
	if err != nil {
		t.Fatalf("GetLinks 报错: %v", err)
	}
	if len(links) != 1 || links[0].Title != "启用友链" {
		t.Errorf("前台友链 = %+v, want 仅启用友链", links)
	}
}

func TestGetLinkList_AllOrderedByTimeDesc(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewLinkService(testutil.NewTestDB(t))

	links, err := s.GetLinkList()
	if err != nil {
		t.Fatalf("GetLinkList 报错: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("len = %d, want 2", len(links))
	}
	if links[0].Title != "启用友链" {
		t.Errorf("应按时间倒序, 首条 = %q", links[0].Title)
	}
}

func TestLinkCRUD(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewLinkService(testutil.NewTestDB(t))

	created := model.FriendLink{Title: "新友链", LinkUrl: "https://new.example.com", State: 1, Type: 1}
	if err := s.Save(created); err != nil {
		t.Fatalf("Save 报错: %v", err)
	}
	list, _ := s.GetLinkList()
	var id uint32
	for _, l := range list {
		if l.Title == "新友链" {
			id = l.Id
		}
	}
	if id == 0 {
		t.Fatal("新建友链未找到")
	}

	got, err := s.GetLink(id)
	if err != nil {
		t.Fatalf("GetLink 报错: %v", err)
	}
	if got.LinkUrl != "https://new.example.com" {
		t.Errorf("LinkUrl = %q", got.LinkUrl)
	}

	got.Title = "新友链-改"
	if err := s.Update(*got); err != nil {
		t.Fatalf("Update 报错: %v", err)
	}
	updated, _ := s.GetLink(id)
	if updated.Title != "新友链-改" {
		t.Errorf("Title = %q, want 新友链-改", updated.Title)
	}

	if err := s.Delete(int(id)); err != nil {
		t.Fatalf("Delete 报错: %v", err)
	}
	if _, err := s.GetLink(id); err == nil {
		t.Error("删除后应查不到")
	}
}
