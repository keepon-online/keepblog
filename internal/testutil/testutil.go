// Package testutil 为业务测试提供独立的 SQLite 测试库。
//
// 特征测试（characterization tests）依赖 global.GORM 是导出变量这一现状：
// 测试进程内首次调用 NewTestDB 时建立临时库并挂到 global.GORM，进程存续期间
// 复用同一实例（Go 的每个测试包是独立进程，包内测试串行执行，因此安全）。
// 阶段 3 引入依赖注入后，各服务可直接以 NewTestDB 返回值构造，不再依赖全局。
package testutil

import (
	"os"
	"sync"
	"testing"
	"time"

	"gitee.com/jieepre/keepblog/global"
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/model/system"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// TestUserPassword 是 fixture 管理员的明文密码，哈希值见 testUserPasswordHash
// （预生成，避免每个测试进程都跑一次 cost=14 的 bcrypt）。
const TestUserPassword = "admin-password-123"

// testUserPasswordHash 是 TestUserPassword 的 bcrypt 哈希。
// 故意用 cost 10（生产为 14）：校验耗时由存储哈希的 cost 决定，
// 降低 cost 让登录相关测试在 -race 下不至于每个用例耗时数秒。
const testUserPasswordHash = "$2a$10$y.oJlTTn4DC9tpHTdht8OeXR4VTeeFE6wGbnCyvMVGaxW061C1jxC"

var (
	dbOnce sync.Once
	testDB *gorm.DB
	dbErr  error
)

// NewTestDB 返回挂载好 fixture 的测试数据库（进程内单例），
// 并保证 global.GORM 指向它。global.GORM 只在首次调用时赋值：
// 服务里有 fire-and-forget 的 goroutine（如异步登录日志）可能在测试
// 结束后仍读取该全局变量，重复赋值会与这些读取构成数据竞争。
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dbOnce.Do(func() {
		dir, err := os.MkdirTemp("", "keepblog-test-*")
		if err != nil {
			dbErr = err
			return
		}
		db, err := gorm.Open(sqlite.Open("file:"+dir+"/test.db?mode=rwc"), &gorm.Config{
			Logger: logger.Discard,
			// 与 internal/core.InitDB 保持一致：表名使用单数形式
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		})
		if err != nil {
			dbErr = err
			return
		}
		if err := db.AutoMigrate(
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
		); err != nil {
			dbErr = err
			return
		}
		if err := seed(db); err != nil {
			dbErr = err
			return
		}
		testDB = db
		global.GORM = db
	})
	if dbErr != nil {
		t.Fatalf("初始化测试数据库失败: %v", dbErr)
	}

	return testDB
}

// ts 是 fixture 用的固定时间锚点，保证测试断言可复现。
func ts(day string) uint64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", day, time.Local)
	if err != nil {
		panic(err)
	}
	return uint64(t.Unix())
}

// seed 写入确定性 fixture。时间线（pub_time = create_time）：
//
//	2024-01-10 P1(go, 置顶)   2024-02-10 P2(mysql)
//	2024-03-10 P3(go)          2024-03-15 P4(sqlite)
//	2024-04-01 P5(mysql)       2024-05-01..07 P6..P12 填充（cat1，无标签）
//	P90 草稿（未发布）         P91 已删除
func seed(db *gorm.DB) error {
	categories := []model.Category{
		{CategoryId: 1, CategoryName: "Go", Note: "Go 语言", State: 1, CreateTime: ts("2024-01-01 00:00:00")},
		{CategoryId: 2, CategoryName: "MySQL", Note: "数据库", State: 1, CreateTime: ts("2024-01-01 00:00:01")},
		{CategoryId: 3, CategoryName: "隐藏分类", Note: "停用", State: 0, CreateTime: ts("2024-01-01 00:00:02")},
	}
	tags := []model.Tag{
		{TagId: 1, TagName: "go", CreateTime: ts("2024-01-01 00:00:00")},
		{TagId: 2, TagName: "mysql", CreateTime: ts("2024-01-01 00:00:01")},
		{TagId: 3, TagName: "sqlite", CreateTime: ts("2024-01-01 00:00:02")},
	}

	newPost := func(id uint64, title, slug, summary, content string, cat uint32, top uint8, pub string, words uint32) model.Post {
		t := ts(pub)
		return model.Post{
			PostId: id, Title: title, PostSlug: slug, Author: "测试作者",
			CoverImage: "cover-" + slug + ".jpg", PostContent: content,
			Summary: summary, Type: 1, Top: top, WordCount: words,
			IsPublished: 1, IsDeleted: 0, Status: 1,
			CategoryId: cat, CreateTime: t, PubTime: t, LastModifiedTime: t,
		}
	}

	posts := []model.Post{
		newPost(1, "Go 快速入门", "go-quickstart", "Go 语言入门摘要", "Go 是一门编译型语言，本文介绍 go 基础。", 1, 1, "2024-01-10 10:00:00", 100),
		newPost(2, "MySQL 统计实战", "mysql-stats", "MySQL 统计与优化摘要", "mysql 的 group by 统计示例。", 2, 0, "2024-02-10 10:00:00", 200),
		newPost(3, "Go 并发模式", "go-concurrency", "goroutine 与 channel 摘要", "Go 并发编程模式。", 1, 0, "2024-03-10 10:00:00", 50),
		newPost(4, "SQLite 迁移指南", "sqlite-migrate", "SQLite 迁移摘要", "SQLite 数据迁移。", 1, 0, "2024-03-15 10:00:00", 40),
		newPost(5, "Redis 缓存实践", "redis-cache", "Redis 缓存摘要", "Redis 实践。", 2, 0, "2024-04-01 10:00:00", 60),
	}
	for i := 6; i <= 12; i++ {
		day := time.Date(2024, 5, i-5, 10, 0, 0, 0, time.Local)
		pub := day.Format("2006-01-02 15:04:05")
		posts = append(posts, newPost(uint64(i), "填充文章"+string(rune('A'+i-6)), "filler-"+string(rune('a'+i-6)),
			"填充摘要", "填充内容。", 1, 0, pub, 10))
	}
	draft := newPost(90, "未发布草稿", "draft-post", "草稿摘要", "草稿内容。", 1, 0, "2024-06-01 10:00:00", 0)
	draft.IsPublished = 0
	deleted := newPost(91, "已删除文章", "deleted-post", "删除摘要", "删除内容。", 1, 0, "2024-06-02 10:00:00", 0)
	deleted.IsDeleted = 1
	posts = append(posts, draft, deleted)

	postTags := []model.PostTag{
		{PostId: 1, TagId: 1},
		{PostId: 2, TagId: 2},
		{PostId: 3, TagId: 1},
		{PostId: 4, TagId: 3},
		{PostId: 5, TagId: 2},
		{PostId: 91, TagId: 1}, // 已删除文章的标签：前台标签云应排除
	}

	links := []model.FriendLink{
		{Id: 1, Title: "启用友链", LinkUrl: "https://a.example.com", State: 1, Type: 1, CreateTime: ts("2024-02-01 00:00:00")},
		{Id: 2, Title: "停用友链", LinkUrl: "https://b.example.com", State: 0, Type: 1, CreateTime: ts("2024-01-01 00:00:00")},
	}

	musics := []model.Music{
		{Id: 1, Name: "启用曲B", Artist: "artist", Url: "https://cdn.example.com/b.mp3", Sort: 2, State: 1, CreateTime: 1},
		{Id: 2, Name: "启用曲A", Artist: "artist", Url: "https://cdn.example.com/a.mp3", Sort: 1, State: 1, CreateTime: 2},
		{Id: 3, Name: "停用曲", Artist: "artist", Url: "https://cdn.example.com/c.mp3", Sort: 3, State: 0, CreateTime: 3},
	}

	user := model.User{
		UserId: 1, Username: "admin", Password: testUserPasswordHash,
		NickName: "管理员", Email: "admin@example.com", Sex: 1,
	}

	website := system.WebSite{
		Id: 1, Title: "测试站点", URL: "https://www.example.com",
		SiteStartDate: "2024-01-01", CreatedAt: int(ts("2024-01-01 00:00:00")),
	}

	for _, batch := range []any{categories, tags, posts, postTags, links, musics} {
		if err := db.Create(batch).Error; err != nil {
			return err
		}
	}
	for _, one := range []any{
		&model.About{Id: 1, Title: "关于", Note: "# 关于博客", CreateTime: ts("2024-01-01 00:00:00")},
		&user, &website,
	} {
		if err := db.Create(one).Error; err != nil {
			return err
		}
	}
	// GORM 对带 default 标签的零值字段会跳过赋值（State:0 落库变默认值 1），
	// 停用曲的状态需在创建后强制写回
	if err := db.Model(&model.Music{}).Where("id = ?", 3).UpdateColumn("state", 0).Error; err != nil {
		return err
	}
	return nil
}
