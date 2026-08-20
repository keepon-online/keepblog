package global

import (
	"gorm.io/gorm"
)

// GORM 全局数据库句柄。
//
// 依赖注入改造完成前，service 层仍从这里取连接（历史现状，约 110 处引用）；
// 阶段 3 将由 AppService 构造时显式传入 *gorm.DB，届时本包只剩
// app.Initialize 里的一处赋值，直至最终删除。
// 原先并存于此的用户级 token 黑名单已合并到 pkg/jwttoken（user_blacklist.go）。
var GORM *gorm.DB
