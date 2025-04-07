package core

import (
	"fmt"
	"github.com/gookit/slog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitDB() *gorm.DB {
	const (
		maxOpenConns    = 10
		maxIdleConns    = 5 // 调整为更合理的空闲连接数
		connMaxLifetime = 60 * time.Minute
	)

	newLogger := logger.New(
		Writer{}, // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	dsn := "file:site.db?&mode=rwc"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger, // 日志配置
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
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

	// 使用插件
	//db.Use(nil)
	return db
}

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	slog.Infof(format, args...)
}
