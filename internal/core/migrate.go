package core

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
)

// InitResource 初始化数据库种子数据。
// schema 和索引由 migrations.Run 负责；数据库连接由调用方先行建立。
func InitResource() {
	if !hasData() {
		initAdmin()
		initData()
	}
	// 默认音乐独立判断：仅当 music 表为空时写入，不依赖 user 表是否已有数据，
	// 这样已初始化的旧库升级后也能补上默认音乐。
	initDefaultMusic()
}

func hasData() bool {
	user := model.User{}
	global.GORM.Model(user).First(&user)
	return user.UserId != 0
}
