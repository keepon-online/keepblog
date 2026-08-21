package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitee.com/jieepre/go-site/config"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
)

var (
	rdb    *redis.Client
	ctx    = context.Background()
	Enable = false
)

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg := config.Get().Redis
	if cfg == nil || !cfg.Enable {
		slog.Info("Redis is disabled")
		return nil
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.Database,
		PoolSize:     100,             // 连接池大小
		MinIdleConns: 10,              // 最小空闲连接
		MaxRetries:   3,               // 最大重试次数
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读超时
		WriteTimeout: 3 * time.Second, // 写超时
	})

	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		slog.Errorf("Failed to connect to Redis: %v", err)
		return err
	}

	Enable = true
	slog.Info("Redis connected successfully")
	return nil
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	if !Enable {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func Get(key string, dest interface{}) error {
	if !Enable {
		return redis.Nil
	}

	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Exists 检查键是否存在
func Exists(key string) bool {
	if !Enable {
		return false
	}

	count, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return count > 0
}

// Delete 删除缓存
func Delete(key string) error {
	if !Enable {
		return nil
	}

	return rdb.Del(ctx, key).Err()
}

// DeletePattern 根据模式删除键（使用 SCAN 避免阻塞）
func DeletePattern(pattern string) error {
	if !Enable {
		return nil
	}

	var cursor uint64
	var deletedCount int64

	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			deleted, err := rdb.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			deletedCount += deleted
		}

		if cursor == 0 {
			break
		}
	}

	if deletedCount > 0 {
		slog.Debugf("Deleted %d keys matching pattern: %s", deletedCount, pattern)
	}

	return nil
}

// Incr 递增
func Incr(key string) (int64, error) {
	if !Enable {
		return 0, nil
	}

	return rdb.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	if !Enable {
		return nil
	}

	return rdb.Expire(ctx, key, expiration).Err()
}

// Ping 健康检查
func Ping() error {
	if !Enable || rdb == nil {
		return nil
	}
	return rdb.Ping(ctx).Err()
}

// Close 关闭Redis连接
func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
