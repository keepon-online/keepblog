package middleware

import (
	"sync"
	"time"
)

// windowLimiter 内存降级用的固定窗口计数器：语义与 Redis 路径（Incr + Expire）
// 保持一致——按「max 次 / window 窗口」计数，窗口到点自动清零。
type windowLimiter struct {
	mu      sync.Mutex
	windows map[string]*windowEntry
}

type windowEntry struct {
	count   int
	resetAt time.Time
}

var globalWindowLimiter = &windowLimiter{windows: make(map[string]*windowEntry)}

// Allow 判断 key 在当前窗口内是否仍可用，并计入本次请求。
func (w *windowLimiter) Allow(key string, max int, window time.Duration) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	e, ok := w.windows[key]
	if !ok || now.After(e.resetAt) {
		e = &windowEntry{resetAt: now.Add(window)}
		w.windows[key] = e
	}
	e.count++

	// 简单清理：条目过多时清掉过期窗口，避免只增不减
	if len(w.windows) > 4096 {
		for k, v := range w.windows {
			if now.After(v.resetAt) {
				delete(w.windows, k)
			}
		}
	}
	return e.count <= max
}