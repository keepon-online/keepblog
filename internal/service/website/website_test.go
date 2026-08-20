package website

import (
	"testing"

	"gitee.com/jieepre/go-site/internal/model/system"
	"gitee.com/jieepre/go-site/internal/testutil"
)

func TestGetWebSite(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewWebSiteService()

	info, err := s.GetWebSite()
	if err != nil {
		t.Fatalf("GetWebSite 报错: %v", err)
	}
	if info.Title != "测试站点" {
		t.Errorf("Title = %q, want 测试站点", info.Title)
	}
	if info.SiteStartDate != "2024-01-01" {
		t.Errorf("SiteStartDate = %q", info.SiteStartDate)
	}
}

func TestSave(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewWebSiteService()

	updated := system.WebSite{Id: 1, Title: "改后站点", URL: "https://update.example.com"}
	if err := s.Save(updated); err != nil {
		t.Fatalf("Save 报错: %v", err)
	}
	info, _ := s.GetWebSite()
	if info.Title != "改后站点" {
		t.Errorf("Title = %q, want 改后站点", info.Title)
	}
}
