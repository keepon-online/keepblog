package core

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/system"
	"github.com/gookit/slog"
)

// InitResource 初始化数据库资源：建表 → 首次启动种子数据 → 索引。
// 数据库连接由调用方（app.Initialize）先行建立并挂到 global.GORM。
func InitResource() {
	initTable()
	if !hasData() {
		initAdmin()
		initData()
	}
	// 默认音乐独立判断：仅当 music 表为空时写入，不依赖 user 表是否已有数据，
	// 这样已初始化的旧库升级后也能补上默认音乐。
	initDefaultMusic()
	// 初始化数据库索引
	InitIndexes()
}

func hasData() bool {
	user := model.User{}
	global.GORM.Model(user).First(&user)
	return user.UserId != 0
}

func initTable() {
	err := global.GORM.AutoMigrate(
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
	if err != nil {
		slog.Errorf("初始化Table失败:%s", err.Error())
	}
	slog.Info("初始化Table")
}
