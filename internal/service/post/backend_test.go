package post

import (
	"testing"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/testutil"
)

// 保存/更新文章时 word_count 由服务端按正文计算，调用方传入的值不生效。
// 用未发布文章测试，避免影响特征测试对已发布文章数的断言。

func TestSavePost_ComputesWordCount(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewPostService(db)

	postObj := model.Post{
		Title:       "字数统计测试",
		PostContent: "# 你好世界\n\nhello world",
		WordCount:   999, // 调用方传入值应被覆盖
		IsPublished: 0, IsDeleted: 0,
	}
	id, err := s.SavePost(postObj)
	if err != nil {
		t.Fatalf("SavePost 报错: %v", err)
	}

	var saved model.Post
	if err := db.Where("post_id = ?", id).First(&saved).Error; err != nil {
		t.Fatalf("查询保存结果报错: %v", err)
	}
	if want := uint32(4 + 10); saved.WordCount != want {
		t.Errorf("WordCount = %d, want %d（你好世界=4 + hello world=10）", saved.WordCount, want)
	}
}

func TestUpdatePost_RecalculatesWordCount(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewPostService(db)

	// 直接构造一篇 word_count 为 0 的存量文章，模拟未回填的旧数据。
	old := model.Post{
		PostId: 8001, Title: "旧文章", PostContent: "旧内容",
		WordCount: 0, IsPublished: 0, IsDeleted: 0,
	}
	if err := db.Create(&old).Error; err != nil {
		t.Fatalf("构造存量文章报错: %v", err)
	}

	old.PostContent = "全新的中文内容加 english words"
	if err := s.UpdatePost(old); err != nil {
		t.Fatalf("UpdatePost 报错: %v", err)
	}

	var saved model.Post
	if err := db.Where("post_id = ?", 8001).First(&saved).Error; err != nil {
		t.Fatalf("查询更新结果报错: %v", err)
	}
	if want := uint32(7 + 13); saved.WordCount != want {
		t.Errorf("WordCount = %d, want %d（全新的中文内容=7 + english words=13）", saved.WordCount, want)
	}
}
