package querybuilder

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type builderFixture struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func (builderFixture) TableName() string { return "builder_fixture" }

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&builderFixture{}); err != nil {
		t.Fatalf("migrate fixture: %v", err)
	}
	if err := db.Create(&builderFixture{Name: "first"}).Error; err != nil {
		t.Fatalf("seed fixture: %v", err)
	}
	return db
}

func TestBuilderCRUDChain(t *testing.T) {
	db := testDB(t)
	qb := New(db, "builder_fixture").Select("name").Order("id DESC")

	var rows []builderFixture
	if err := qb.Find(&rows); err != nil {
		t.Fatalf("Find: %v", err)
	}
	if len(rows) != 1 || rows[0].Name != "first" {
		t.Fatalf("rows = %+v", rows)
	}

	var count int64
	if err := New(db, "builder_fixture").Count(&count); err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}

	var one builderFixture
	if err := New(db, "builder_fixture").First(&one); err != nil {
		t.Fatalf("First: %v", err)
	}
	if one.Name != "first" {
		t.Errorf("one.Name = %q", one.Name)
	}
}

func TestBuilderScopesAndPagination(t *testing.T) {
	db := testDB(t)
	qb := New(db, "builder_fixture").
		Scope(func(db *gorm.DB) *gorm.DB { return db.Where("name = ?", "first") }).
		Limit(1).
		Offset(0)

	var rows []builderFixture
	if err := qb.Find(&rows); err != nil {
		t.Fatalf("Find: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("rows = %d, want 1", len(rows))
	}
}
