package sidebar

import (
	"testing"

	"gitee.com/jieepre/go-site/internal/testutil"
)

// Sidebar 聚合了各服务的前台数据，钉住聚合结构与计数。
func TestSidebar_Aggregate(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewSidebarService(testutil.NewTestDB(t))

	sb := s.Sidebar()

	// 12 篇已发布、2 个分类、3 个标签
	if sb.CardInfo.Post != 12 {
		t.Errorf("CardInfo.Post = %d, want 12", sb.CardInfo.Post)
	}
	if sb.CardInfo.Category != 2 {
		t.Errorf("CardInfo.Category = %d, want 2", sb.CardInfo.Category)
	}
	if sb.CardInfo.Tag != 3 {
		t.Errorf("CardInfo.Tag = %d, want 3", sb.CardInfo.Tag)
	}

	if len(sb.LatestPosts) != 5 {
		t.Errorf("LatestPosts len = %d, want 5", len(sb.LatestPosts))
	}

	// 分类计数同 category.GetCategories：无文章的分类不出现
	if len(sb.Category) != 2 {
		t.Errorf("Category len = %d, want 2", len(sb.Category))
	}

	// 按年月分组：5 组，2024/05 含 7 篇
	if len(sb.SidebarArchives) != 5 {
		t.Fatalf("SidebarArchives len = %d, want 5", len(sb.SidebarArchives))
	}
	for _, a := range sb.SidebarArchives {
		if a.Year == "2024/05" && a.Total != 7 {
			t.Errorf("2024/05 = %d 篇, want 7", a.Total)
		}
	}

	// 站点资讯：字数合计 100+200+50+40+60+7*10=520
	if sb.WebInfo.PostCount != 12 {
		t.Errorf("WebInfo.PostCount = %d, want 12", sb.WebInfo.PostCount)
	}
	if sb.WebInfo.TotalWordCount != 520 {
		t.Errorf("WebInfo.TotalWordCount = %d, want 520", sb.WebInfo.TotalWordCount)
	}
	if sb.WebInfo.SiteStartDate != "2024-01-01" {
		t.Errorf("WebInfo.SiteStartDate = %q, want 2024-01-01", sb.WebInfo.SiteStartDate)
	}
	if sb.WebInfo.LastUpdateTime == "" {
		t.Error("LastUpdateTime 不应为空")
	}
}
