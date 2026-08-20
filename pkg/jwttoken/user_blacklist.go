package jwttoken

import "sync"

// 用户级令牌失效名单：密码修改后使该用户名下既有令牌全部失效。
// 与按 token 串的黑名单（blacklist.go）互补：那份针对单个 token 注销，
// 这份按用户名 + 签发时间批量失效。原先散落在 global 包，统一收拢于此。
var (
	userBlacklist   = make(map[string]int64)
	userBlacklistMu sync.RWMutex
)

// InvalidateUserTokens 使用户的所有令牌失效（在密码修改时调用）。
// timestamp 之后判定：签发时间早于该时间戳的令牌均无效。
func InvalidateUserTokens(username string, timestamp int64) {
	userBlacklistMu.Lock()
	defer userBlacklistMu.Unlock()
	userBlacklist[username] = timestamp
}

// IsUserTokenValid 检查令牌是否未被用户级撤销
func IsUserTokenValid(username string, issuedAt int64) bool {
	userBlacklistMu.RLock()
	defer userBlacklistMu.RUnlock()
	invalidatedAt, exists := userBlacklist[username]
	if !exists {
		return true
	}
	return issuedAt > invalidatedAt
}
