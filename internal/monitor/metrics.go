package monitor

import (
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics 性能指标
type Metrics struct {
	// HTTP指标
	RequestsTotal   int64            `json:"requests_total"`
	RequestDuration time.Duration    `json:"avg_request_duration"`
	StatusCodes     map[string]int64 `json:"status_codes"`

	// 系统指标
	CPUUsage       float64 `json:"cpu_usage_percent"`
	MemoryUsage    int64   `json:"memory_usage_bytes"`
	GoroutineCount int     `json:"goroutine_count"`
	GCCount        uint32  `json:"gc_count"`

	// 自定义指标
	ActiveConnections int     `json:"active_connections"`
	CacheHitRate      float64 `json:"cache_hit_rate"`

	mutex         sync.RWMutex
	startTime     time.Time
	totalDuration time.Duration
	cacheHits     int64
	cacheMisses   int64
}

// NewMetrics 创建指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		StatusCodes: make(map[string]int64),
		startTime:   time.Now(),
	}
}

// RecordRequest 记录请求指标
func (m *Metrics) RecordRequest(statusCode int, duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.RequestsTotal++
	m.totalDuration += duration
	m.RequestDuration = m.totalDuration / time.Duration(m.RequestsTotal)

	statusKey := getStatusGroup(statusCode)
	m.StatusCodes[statusKey]++
}

// RecordCacheHit 记录缓存命中
func (m *Metrics) RecordCacheHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.cacheHits++
	m.updateCacheHitRate()
}

// RecordCacheMiss 记录缓存未命中
func (m *Metrics) RecordCacheMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.cacheMisses++
	m.updateCacheHitRate()
}

// updateCacheHitRate 更新缓存命中率
func (m *Metrics) updateCacheHitRate() {
	total := m.cacheHits + m.cacheMisses
	if total > 0 {
		m.CacheHitRate = float64(m.cacheHits) / float64(total) * 100
	}
}

// UpdateSystemMetrics 更新系统指标
func (m *Metrics) UpdateSystemMetrics() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.MemoryUsage = int64(memStats.Alloc)
	m.GoroutineCount = runtime.NumGoroutine()
	m.GCCount = memStats.NumGC
}

// GetMetrics 获取当前指标
func (m *Metrics) GetMetrics() Metrics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 更新系统指标
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 创建副本返回
	statusCodes := make(map[string]int64)
	for k, v := range m.StatusCodes {
		statusCodes[k] = v
	}

	return Metrics{
		RequestsTotal:     m.RequestsTotal,
		RequestDuration:   m.RequestDuration,
		StatusCodes:       statusCodes,
		MemoryUsage:       int64(memStats.Alloc),
		GoroutineCount:    runtime.NumGoroutine(),
		GCCount:           memStats.NumGC,
		ActiveConnections: m.ActiveConnections,
		CacheHitRate:      m.CacheHitRate,
	}
}

// MetricsHandler HTTP指标处理器
func (m *Metrics) MetricsHandler(c *gin.Context) {
	metrics := m.GetMetrics()
	c.JSON(200, metrics)
}

// getStatusGroup 获取状态码分组
func getStatusGroup(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}

// MetricsMiddleware 指标收集中间件
func MetricsMiddleware(metrics *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录指标
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		metrics.RecordRequest(statusCode, duration)
	}
}
