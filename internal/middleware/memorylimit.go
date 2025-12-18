package middleware

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// MemoryRateLimiter 内存限流器，用于 Redis 不可用时的降级方案
type MemoryRateLimiter struct {
	limiters sync.Map
	mu       sync.RWMutex
	limit    rate.Limit
	burst    int
}

// limiterEntry 限流器条目（带过期时间）
type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// 全局内存限流器实例
var (
	globalMemoryLimiter *MemoryRateLimiter
	memoryLimiterOnce   sync.Once
)

// GetMemoryLimiter 获取全局内存限流器实例
func GetMemoryLimiter() *MemoryRateLimiter {
	memoryLimiterOnce.Do(func() {
		globalMemoryLimiter = NewMemoryRateLimiter(1, 60) // 每秒1个请求，突发60个
		go globalMemoryLimiter.cleanupLoop()
	})
	return globalMemoryLimiter
}

// NewMemoryRateLimiter 创建新的内存限流器
// ratePerSecond: 每秒请求数
// burst: 突发请求数
func NewMemoryRateLimiter(ratePerSecond float64, burst int) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		limit: rate.Limit(ratePerSecond),
		burst: burst,
	}
}

// Allow 检查是否允许请求
func (m *MemoryRateLimiter) Allow(key string) bool {
	entry, _ := m.limiters.LoadOrStore(key, &limiterEntry{
		limiter:  rate.NewLimiter(m.limit, m.burst),
		lastSeen: time.Now(),
	})

	le := entry.(*limiterEntry)
	le.lastSeen = time.Now()

	return le.limiter.Allow()
}

// AllowN 检查是否允许 n 个请求
func (m *MemoryRateLimiter) AllowN(key string, n int) bool {
	entry, _ := m.limiters.LoadOrStore(key, &limiterEntry{
		limiter:  rate.NewLimiter(m.limit, m.burst),
		lastSeen: time.Now(),
	})

	le := entry.(*limiterEntry)
	le.lastSeen = time.Now()

	return le.limiter.AllowN(time.Now(), n)
}

// SetRate 设置限流速率
func (m *MemoryRateLimiter) SetRate(ratePerSecond float64, burst int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limit = rate.Limit(ratePerSecond)
	m.burst = burst
}

// cleanupLoop 定期清理过期的限流器
func (m *MemoryRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.cleanup()
	}
}

// cleanup 清理超过10分钟未使用的限流器
func (m *MemoryRateLimiter) cleanup() {
	threshold := time.Now().Add(-10 * time.Minute)

	m.limiters.Range(func(key, value interface{}) bool {
		le := value.(*limiterEntry)
		if le.lastSeen.Before(threshold) {
			m.limiters.Delete(key)
		}
		return true
	})
}

// Reset 重置指定 key 的限流器
func (m *MemoryRateLimiter) Reset(key string) {
	m.limiters.Delete(key)
}

// Stats 返回当前活跃的限流器数量
func (m *MemoryRateLimiter) Stats() int {
	count := 0
	m.limiters.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
