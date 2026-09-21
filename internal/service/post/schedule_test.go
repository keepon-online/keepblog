package post

import (
	"testing"
	"time"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/testutil"
)

// TestScheduledPost_InvisibleUntilDue 定时发布核心语义：
// is_published=1 且 pub_time 在未来的文章，前台所有路径（详情/最新/归档/
// 系列）都不可见；pub_time 过去后自然可见，无需任何翻转任务。
func TestScheduledPost_InvisibleUntilDue(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewPostService(db)

	future := uint64(time.Now().Add(1 * time.Hour).Unix())
	past := uint64(time.Now().Add(-1 * time.Hour).Unix())

	// 一篇定时、一篇同系列已到点，用于同时验证系列查询
	scheduled := model.Post{
		Title: "定时文章", IsPublished: 1, Series: "系列A",
		PubTime: future, CreateTime: uint64(time.Now().Unix()),
	}
	if err := db.Create(&scheduled).Error; err != nil {
		t.Fatal(err)
	}
	due := model.Post{
		Title: "已到点文章", IsPublished: 1, Series: "系列A",
		PubTime: past, CreateTime: uint64(time.Now().Unix()),
	}
	if err := db.Create(&due).Error; err != nil {
		t.Fatal(err)
	}

	// 详情：定时文章零值（不可见），已到点文章可见
	got, err := s.GetPost(int(scheduled.PostId))
	if err != nil || got.PostId != 0 {
		t.Errorf("GetPost(定时) = (%d, %v), want (0, nil)", got.PostId, err)
	}
	got, err = s.GetPost(int(due.PostId))
	if err != nil || got.PostId != due.PostId {
		t.Errorf("GetPost(已到点) = (%d, %v), want (%d, nil)", got.PostId, err, due.PostId)
	}

	// 最新文章：定时文章不出现
	latest, err := s.GetLatestPosts()
	if err != nil {
		t.Fatal(err)
	}
	for _, lp := range latest {
		if lp.Title == "定时文章" {
			t.Error("GetLatestPosts 泄露了未到点的定时文章")
		}
	}

	// 系列：只包含已到点文章
	series, err := s.GetSeriesPosts("系列A")
	if err != nil {
		t.Fatal(err)
	}
	if len(series) != 1 || series[0].Title != "已到点文章" {
		t.Errorf("GetSeriesPosts = %+v, want 仅[已到点文章]", series)
	}

	// 搜索：定时文章不出现
	results, err := s.Search("定时文章")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Title == "定时文章" {
			t.Error("Search 泄露了未到点的定时文章")
		}
	}
}

// TestUpdatePost_PreservesPubTime 回归测试：pub_time 不再随编辑被重置。
// 旧行为（gorm autoUpdateTime）会在每次 Save 时把发布时间覆盖为当前时刻。
func TestUpdatePost_PreservesPubTime(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewPostService(db)

	past := uint64(time.Now().Add(-48 * time.Hour).Unix())
	p := model.Post{
		Title: "旧文章", IsPublished: 1, PubTime: past,
		CreateTime: past, PostContent: "内容",
	}
	if err := db.Create(&p).Error; err != nil {
		t.Fatal(err)
	}

	// 编辑时不带 pubTime（前端未开启定时发布的默认形态）
	if err := s.UpdatePost(model.Post{
		PostId: p.PostId, Title: "旧文章改", PostContent: "新内容",
		CreateTime: past, IsPublished: 1,
	}); err != nil {
		t.Fatal(err)
	}

	var got model.Post
	db.First(&got, p.PostId)
	if got.PubTime != past {
		t.Errorf("pub_time = %d, want 保持原值 %d（编辑不得重置发布时间）", got.PubTime, past)
	}

	// 显式指定未来时间（定时发布）则原样生效
	future := uint64(time.Now().Add(24 * time.Hour).Unix())
	if err := s.UpdatePost(model.Post{
		PostId: p.PostId, Title: "旧文章改", PostContent: "新内容",
		CreateTime: past, IsPublished: 1, PubTime: future,
	}); err != nil {
		t.Fatal(err)
	}
	db.First(&got, p.PostId)
	if got.PubTime != future {
		t.Errorf("pub_time = %d, want 定时时间 %d", got.PubTime, future)
	}
}

// TestSavePost_DefaultPubTimeNow 新建文章未指定发布时间时默认当前时间。
func TestSavePost_DefaultPubTimeNow(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewPostService(db)

	before := time.Now().Unix()
	id, err := s.SavePost(model.Post{Title: "新文章", PostContent: "内容"})
	if err != nil {
		t.Fatal(err)
	}
	var got model.Post
	db.First(&got, id)
	if got.PubTime < uint64(before) || got.PubTime > uint64(time.Now().Unix()) {
		t.Errorf("pub_time = %d, want 创建时刻附近", got.PubTime)
	}
}

// TestGetSeriesPosts_OrderAscExcludeOthers 系列查询按发布时间升序，
// 排除草稿、已删除与其他系列文章。
func TestGetSeriesPosts_OrderAscExcludeOthers(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewPostService(db)

	base := uint64(time.Now().Add(-72 * time.Hour).Unix())
	mk := func(title, series string, pub uint64, published uint8) model.Post {
		return model.Post{Title: title, Series: series, IsPublished: published,
			PubTime: pub, CreateTime: pub, IsDeleted: 0}
	}
	posts := []model.Post{
		mk("第二篇", "系列B", base+7200, 1),
		mk("第一篇", "系列B", base+3600, 1),
		mk("第三篇", "系列B", base+10800, 1),
		mk("草稿篇", "系列B", base+5400, 0),
		mk("别系列", "系列C", base+6300, 1),
		mk("删除篇", "系列B", base+9000, 1),
	}
	for i := range posts {
		if err := db.Create(&posts[i]).Error; err != nil {
			t.Fatal(err)
		}
	}
	db.Model(&model.Post{}).Where("title = ?", "删除篇").Update("is_deleted", 1)

	got, err := s.GetSeriesPosts("系列B")
	if err != nil {
		t.Fatal(err)
	}
	wantTitles := []string{"第一篇", "第二篇", "第三篇"}
	if len(got) != len(wantTitles) {
		t.Fatalf("GetSeriesPosts 返回 %d 篇 %+v, want %d 篇", len(got), got, len(wantTitles))
	}
	for i, w := range wantTitles {
		if got[i].Title != w {
			t.Errorf("系列第 %d 篇 = %q, want %q（应按发布时间升序）", i+1, got[i].Title, w)
		}
	}
}
