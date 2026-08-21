package migrations

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openMigrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return db
}

func TestRunCreatesSchemaAndIsIdempotent(t *testing.T) {
	db := openMigrationDB(t)
	if err := Run(db); err != nil {
		t.Fatalf("first migration run: %v", err)
	}
	if err := Run(db); err != nil {
		t.Fatalf("second migration run: %v", err)
	}

	var count int
	if err := db.Raw("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name IN ('post', 'category', 'system_notice')").Scan(&count).Error; err != nil {
		t.Fatalf("inspect schema: %v", err)
	}
	if count != 3 {
		t.Fatalf("created table count = %d, want 3", count)
	}

	var version int64
	if err := db.Raw("SELECT version_id FROM goose_db_version WHERE is_applied = 1 ORDER BY version_id DESC LIMIT 1").Scan(&version).Error; err != nil {
		t.Fatalf("inspect goose version: %v", err)
	}
	if version != 20260821000100 {
		t.Errorf("migration version = %d, want 20260821000100", version)
	}
}

func TestRunPreservesExistingData(t *testing.T) {
	db := openMigrationDB(t)
	// 模拟旧版本已存在的文章表；其字段覆盖历史索引所需列。
	if err := db.Exec(`CREATE TABLE post (
		post_id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT, post_slug TEXT, author TEXT, cover_image TEXT,
		post_content TEXT, post_content_html TEXT, summary TEXT,
		type INTEGER, top INTEGER, read_count INTEGER, word_count INTEGER,
		is_published INTEGER, is_deleted INTEGER, status INTEGER,
		create_time INTEGER, pub_time INTEGER, last_modified_time INTEGER,
		category_id INTEGER
	)`).Error; err != nil {
		t.Fatalf("create legacy table: %v", err)
	}
	if err := db.Exec("INSERT INTO post(post_id, title) VALUES (7, 'legacy')").Error; err != nil {
		t.Fatalf("insert legacy row: %v", err)
	}

	if err := Run(db); err != nil {
		t.Fatalf("migration run: %v", err)
	}
	var title string
	if err := db.Raw("SELECT title FROM post WHERE post_id = 7").Scan(&title).Error; err != nil {
		t.Fatalf("read legacy row: %v", err)
	}
	if title != "legacy" {
		t.Errorf("legacy title = %q, want legacy", title)
	}
}
