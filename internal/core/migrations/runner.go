package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/model/system"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed sql/*.sql
var migrationFS embed.FS

// Run 将 SQLite schema 通过版本化迁移推进到最新版本。
// 迁移使用 IF NOT EXISTS，因此对已有 AutoMigrate 数据库是无损的；
// 后续字段变更应新增 migration 文件，不再直接修改历史基线。
func Run(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("migration database is nil")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql database: %w", err)
	}
	sqlFS, err := fs.Sub(migrationFS, "sql")
	if err != nil {
		return fmt.Errorf("open embedded migrations: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, sqlDB, sqlFS)
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	ctx := context.Background()
	version, err := provider.GetDBVersion(ctx)
	if err != nil {
		return fmt.Errorf("read migration version: %w", err)
	}
	if version == 0 {
		if err := bootstrapLegacySchema(db); err != nil {
			return fmt.Errorf("bootstrap legacy schema: %w", err)
		}
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	if err := ensureSearchVerificationColumns(db); err != nil {
		return fmt.Errorf("ensure search verification columns: %w", err)
	}
	return nil
}

// bootstrapLegacySchema 仅处理尚未接入 goose 的旧部署：AutoMigrate 会补齐
// 基线缺少的字段/表，但不会删除或重建已有数据。完成后基线 SQL 的
// IF NOT EXISTS 语句负责索引与新表，后续启动不再调用此函数。
func bootstrapLegacySchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Post{},
		&model.Category{},
		&model.FriendLink{},
		&model.Tag{},
		&model.PostTag{},
		&model.About{},
		&model.User{},
		&model.Music{},
		&system.AccessLog{},
		&system.WebSite{},
		&system.LoginLog{},
		&system.Notice{},
	)
}

// ensureSearchVerificationColumns 给 web_site 补 Google/Bing 站长平台验证列
// （2026-09 加入）。不走 goose 迁移文件：SQLite 的 ADD COLUMN 没有
// IF NOT EXISTS，而版本 0 引导路径的 AutoMigrate 已按模型建出这两列，
// SQL 迁移在全新库上必然撞重复列。此处幂等补列，已有列时为两次轻量
// pragma 查询，成本可忽略。
func ensureSearchVerificationColumns(db *gorm.DB) error {
	for _, col := range []string{"google_site", "bing_site"} {
		var exists int64
		if err := db.Raw(
			"SELECT COUNT(1) FROM pragma_table_info('web_site') WHERE name = ?", col,
		).Scan(&exists).Error; err != nil {
			return err
		}
		if exists > 0 {
			continue
		}
		if err := db.Exec(
			fmt.Sprintf("ALTER TABLE web_site ADD COLUMN %s TEXT DEFAULT ''", col),
		).Error; err != nil {
			return err
		}
	}
	return nil
}
