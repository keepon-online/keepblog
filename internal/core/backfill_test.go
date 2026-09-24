package core

import (
	"testing"

	"gitee.com/jieepre/keepblog/internal/core/migrations"
	"gitee.com/jieepre/keepblog/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func newBackfillDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
		Logger: logger.Discard,
		// 与 internal/core.InitDB 保持一致：表名使用单数形式。
		// 否则 gorm Create 会写入 AutoMigrate 顺带建出的复数表（posts），
		// 与基线 SQL 的 post 表错位。
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := migrations.Run(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	return db
}

// fetchPost 每次用零值结构体查询：复用带主键的结构体做 First 目标时，
// gorm 会把旧主键拼进查询条件导致误报 record not found。
func fetchPost(t *testing.T, db *gorm.DB, postId uint64) model.Post {
	t.Helper()
	var p model.Post
	if err := db.Where("post_id = ?", postId).First(&p).Error; err != nil {
		t.Fatalf("查询文章 %d 报错: %v", postId, err)
	}
	return p
}

func TestBackfillPostWordCount(t *testing.T) {
	db := newBackfillDB(t)
	seed := []model.Post{
		{PostId: 1, Title: "待回填", PostContent: "你好世界 hello", WordCount: 0},
		{PostId: 2, Title: "已有计数", PostContent: "保持原样", WordCount: 42},
		{PostId: 3, Title: "空正文", PostContent: "", WordCount: 0},
	}
	for i := range seed {
		if err := db.Create(&seed[i]).Error; err != nil {
			t.Fatalf("构造种子文章 %d 报错: %v", seed[i].PostId, err)
		}
	}

	BackfillPostWordCount(db)

	if got := fetchPost(t, db, 1); got.WordCount != 4+5 {
		t.Errorf("文章 1 WordCount = %d, want %d（你好世界=4 + hello=5）", got.WordCount, 4+5)
	}
	if got := fetchPost(t, db, 2); got.WordCount != 42 {
		t.Errorf("已有计数的文章不应被覆盖，WordCount = %d, want 42", got.WordCount)
	}
	if got := fetchPost(t, db, 3); got.WordCount != 0 {
		t.Errorf("空正文文章 WordCount = %d, want 0", got.WordCount)
	}

	// 幂等：再次回填无变化、不报错。
	BackfillPostWordCount(db)
}

func TestBackfillPostWordCount_DoesNotTouchLastModifiedTime(t *testing.T) {
	db := newBackfillDB(t)
	p := model.Post{PostId: 10, Title: "时间锚点", PostContent: "内容", WordCount: 0,
		LastModifiedTime: 1700000000}
	if err := db.Create(&p).Error; err != nil {
		t.Fatalf("构造种子文章报错: %v", err)
	}

	BackfillPostWordCount(db)

	if got := fetchPost(t, db, 10); got.LastModifiedTime != 1700000000 {
		t.Errorf("回填不应改动 last_modified_time = %d, want 1700000000", got.LastModifiedTime)
	}
}

func TestBackfillPostSlug(t *testing.T) {
	db := newBackfillDB(t)
	seed := []model.Post{
		{PostId: 101, Title: "待补齐slug文章", PostSlug: ""},
		{PostId: 102, Title: "已有slug文章", PostSlug: "customSlug123"},
	}
	for i := range seed {
		if err := db.Create(&seed[i]).Error; err != nil {
			t.Fatalf("构造种子文章 %d 报错: %v", seed[i].PostId, err)
		}
	}

	BackfillPostSlug(db)

	if got := fetchPost(t, db, 101); got.PostSlug == "" {
		t.Errorf("文章 101 的 PostSlug 未被正确补齐")
	}
	if got := fetchPost(t, db, 102); got.PostSlug != "customSlug123" {
		t.Errorf("已有 PostSlug 的文章被意外覆盖: %s, want customSlug123", got.PostSlug)
	}

	// 幂等：再次执行无异常
	BackfillPostSlug(db)
}
