package core

import (
	"fmt"
	"time"

	"github.com/gookit/slog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() *gorm.DB {
	const (
		maxOpenConns    = 25               // 增加最大连接数
		maxIdleConns    = 10               // 增加空闲连接数
		connMaxLifetime = 2 * time.Hour    // 增加连接生命周期
		connMaxIdleTime = 30 * time.Minute // 新增：空闲连接最大存活时间
	)

	newLogger := logger.New(
		Writer{}, // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 降低慢查询阈值
			LogLevel:                  logger.Warn,            // 调整为Warn级别
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
			ParameterizedQueries:      true,                   // 启用参数化查询
		},
	)
	dsn := "file:./data/site.db?&mode=rwc"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		// 性能优化配置
		PrepareStmt:                              true, // 缓存预编译语句
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束，提升性能
	})

	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库实例失败: %v", err))
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(maxOpenConns)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// 设置空闲连接最大存活时间
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	// 使用插件
	//db.Use(nil)
	return db
}

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	slog.Infof(format, args...)
}
