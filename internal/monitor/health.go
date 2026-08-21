package monitor

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"gitee.com/jieepre/keepblog/global"
	"gitee.com/jieepre/keepblog/internal/cache"

	"github.com/gin-gonic/gin"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Checks    map[string]Health `json:"checks"`
	System    *SystemInfo       `json:"system,omitempty"`
}

// Health 单个检查项的健康状态
type Health struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	Uptime     time.Duration `json:"uptime"`
	GoVersion  string        `json:"go_version"`
	Goroutines int           `json:"goroutines"`
	MemAllocMB float64       `json:"mem_alloc_mb"`
	MemSysMB   float64       `json:"mem_sys_mb"`
	NumCPU     int           `json:"num_cpu"`
	NumGC      uint32        `json:"num_gc"`
}

// HealthChecker 健康检查器
type HealthChecker struct {
	startTime time.Time
	version   string
	checkers  map[string]func() Health
	mutex     sync.RWMutex
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(version string) *HealthChecker {
	hc := &HealthChecker{
		startTime: time.Now(),
		version:   version,
		checkers:  make(map[string]func() Health),
	}

	// 注册默认检查项
	hc.RegisterChecker("database", hc.checkDatabase)
	hc.RegisterChecker("redis", hc.checkRedis)

	return hc
}

// RegisterChecker 注册检查项
func (hc *HealthChecker) RegisterChecker(name string, checker func() Health) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.checkers[name] = checker
}

// HealthCheck HTTP处理器
func (hc *HealthChecker) HealthCheck(c *gin.Context) {
	status := hc.GetHealth()

	httpStatus := http.StatusOK
	if status.Status != "UP" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, status)
}

// GetHealth 获取健康状态
func (hc *HealthChecker) GetHealth() HealthStatus {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	checks := make(map[string]Health)
	overallStatus := "UP"

	// 执行所有检查项
	for name, checker := range hc.checkers {
		check := checker()
		checks[name] = check

		if check.Status != "UP" {
			overallStatus = "DOWN"
		}
	}

	return HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   hc.version,
		Checks:    checks,
		System:    hc.getSystemInfo(),
	}
}

// checkDatabase 检查数据库连接
func (hc *HealthChecker) checkDatabase() Health {
	start := time.Now()

	if global.GORM == nil {
		return Health{
			Status:  "DOWN",
			Message: "Database connection not initialized",
			Latency: time.Since(start),
		}
	}

	sqlDB, err := global.GORM.DB()
	if err != nil {
		return Health{
			Status:  "DOWN",
			Message: fmt.Sprintf("Failed to get underlying sql.DB: %v", err),
			Latency: time.Since(start),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return Health{
			Status:  "DOWN",
			Message: fmt.Sprintf("Database ping failed: %v", err),
			Latency: time.Since(start),
		}
	}

	return Health{
		Status:  "UP",
		Message: "Database is healthy",
		Latency: time.Since(start),
	}
}

// checkRedis 检查Redis连接
func (hc *HealthChecker) checkRedis() Health {
	start := time.Now()

	if !cache.Enable {
		return Health{
			Status:  "UP",
			Message: "Redis is disabled",
			Latency: time.Since(start),
		}
	}

	// 使用 PING 命令进行健康检查（更轻量）
	if err := cache.Ping(); err != nil {
		return Health{
			Status:  "DOWN",
			Message: fmt.Sprintf("Redis ping failed: %v", err),
			Latency: time.Since(start),
		}
	}

	return Health{
		Status:  "UP",
		Message: "Redis is healthy",
		Latency: time.Since(start),
	}
}

// getSystemInfo 获取系统信息
func (hc *HealthChecker) getSystemInfo() *SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &SystemInfo{
		Uptime:     time.Since(hc.startTime),
		GoVersion:  runtime.Version(),
		Goroutines: runtime.NumGoroutine(),
		MemAllocMB: float64(m.Alloc) / 1024 / 1024,
		MemSysMB:   float64(m.Sys) / 1024 / 1024,
		NumCPU:     runtime.NumCPU(),
		NumGC:      m.NumGC,
	}
}

// ReadinessCheck 就绪检查（简化版健康检查）
func (hc *HealthChecker) ReadinessCheck(c *gin.Context) {
	// 只检查关键服务
	hc.mutex.RLock()
	dbChecker, hasDB := hc.checkers["database"]
	hc.mutex.RUnlock()

	if hasDB {
		dbHealth := dbChecker()
		if dbHealth.Status != "UP" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "NOT_READY",
				"message": "Database is not ready",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "READY",
		"timestamp": time.Now(),
	})
}

// LivenessCheck 存活检查（最简单的检查）
func (hc *HealthChecker) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ALIVE",
		"timestamp": time.Now(),
		"uptime":    time.Since(hc.startTime).String(),
	})
}
