package jwttoken

import (
	"sync"
	"time"

	"gitee.com/jieepre/go-site/internal/cache"
)

// TokenBlacklist Token 黑名单管理
// 用于 Logout 后使 Token 失效

var (
	// 内存黑名单（Redis 不可用时的备选方案）
	memoryBlacklist    = make(map[string]time.Time)
	memoryBlacklistMu  sync.RWMutex
	blacklistCleanupOn sync.Once
)

const (
	blacklistKeyPrefix = "token_blacklist:"
	blacklistTTL       = 24 * time.Hour // Token 黑名单保留时间（应大于 Token 最大有效期）
)

// init 启动清理协程
func init() {
	blacklistCleanupOn.Do(func() {
		go cleanupBlacklist()
	})
}

// InvalidateToken 使 Token 失效（加入黑名单）
func InvalidateToken(tokenStr string) error {
	expiry := time.Now().Add(blacklistTTL)
	key := blacklistKeyPrefix + tokenStr

	// 尝试写入 Redis
	if cache.Enable {
		if err := cache.Set(key, expiry.Unix(), blacklistTTL); err == nil {
			return nil
		}
	}

	// Redis 不可用时，使用内存黑名单
	memoryBlacklistMu.Lock()
	memoryBlacklist[tokenStr] = expiry
	memoryBlacklistMu.Unlock()

	return nil
}

// IsTokenBlacklisted 检查 Token 是否在黑名单中
func IsTokenBlacklisted(tokenStr string) bool {
	key := blacklistKeyPrefix + tokenStr

	// 先检查 Redis
	if cache.Enable {
		if cache.Exists(key) {
			return true
		}
	}

	// 检查内存黑名单
	memoryBlacklistMu.RLock()
	expiry, exists := memoryBlacklist[tokenStr]
	memoryBlacklistMu.RUnlock()

	if exists && time.Now().Before(expiry) {
		return true
	}

	return false
}

// cleanupBlacklist 定期清理过期的内存黑名单条目
func cleanupBlacklist() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		memoryBlacklistMu.Lock()
		now := time.Now()
		for token, expiry := range memoryBlacklist {
			if now.After(expiry) {
				delete(memoryBlacklist, token)
			}
		}
		memoryBlacklistMu.Unlock()
	}
}
