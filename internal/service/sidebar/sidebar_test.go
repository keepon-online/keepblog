package sidebar

import (
	"testing"

	"gitee.com/jieepre/keepblog/internal/model/system"
	"gitee.com/jieepre/keepblog/internal/testutil"
	"gorm.io/gorm"
)

// Sidebar 聚合了各服务的前台数据，钉住聚合结构与计数。
func TestSidebar_Aggregate(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewSidebarService(db)

	// 访问日志：两个 IP、PV 合计 8（1+3 + 2+2），去重 UV = 2
	seedAccessLogs(t, db)

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

	// 访客数/访问量：来自 system_access_log（去重 IP 口径，与 dashboard 一致）
	if sb.WebInfo.SiteUV != 2 {
		t.Errorf("WebInfo.SiteUV = %d, want 2", sb.WebInfo.SiteUV)
	}
	if sb.WebInfo.SitePV != 8 {
		t.Errorf("WebInfo.SitePV = %d, want 8", sb.WebInfo.SitePV)
	}
}

// seedAccessLogs 写入确定性的访问日志 fixture：仅本包使用，
// 避免污染共享测试库中其他包对 system_access_log 的统计口径。
func seedAccessLogs(t *testing.T, db *gorm.DB) {
	t.Helper()
	logs := []system.AccessLog{
		accessLog(3232235777, "/", 1),  // 192.168.1.1
		accessLog(3232235777, "/about", 3),
		accessLog(3232235521, "/", 2),  // 192.168.0.1
		accessLog(3232235521, "/link", 2),
	}
	for i := range logs {
		if err := db.Create(&logs[i]).Error; err != nil {
			t.Fatalf("构造访问日志报错: %v", err)
		}
	}
}

func accessLog(ip uint32, url string, pv int) system.AccessLog {
	p := int(pv)
	u := uint32(ip)
	return system.AccessLog{Ip: &u, URL: url, PV: &p, UV: &p}
}
