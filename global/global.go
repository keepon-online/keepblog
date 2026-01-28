package global

import (
	"sync"

	"gorm.io/gorm"
)

var (
	GORM *gorm.DB

	// TokenBlacklist 用于存储已失效的令牌（密码修改后使旧令牌失效）
	// key: username, value: 令牌失效时间戳（该时间之前签发的令牌都无效）
	TokenBlacklist     = make(map[string]int64)
	TokenBlacklistLock sync.RWMutex
)

// InvalidateUserTokens 使用户的所有令牌失效（在密码修改时调用）
func InvalidateUserTokens(username string, timestamp int64) {
	TokenBlacklistLock.Lock()
	defer TokenBlacklistLock.Unlock()
	TokenBlacklist[username] = timestamp
}

// IsTokenValid 检查令牌是否有效（未被撤销）
func IsTokenValid(username string, issuedAt int64) bool {
	TokenBlacklistLock.RLock()
	defer TokenBlacklistLock.RUnlock()
	invalidatedAt, exists := TokenBlacklist[username]
	if !exists {
		return true
	}
	// 如果令牌签发时间早于失效时间，则令牌无效
	return issuedAt > invalidatedAt
}
