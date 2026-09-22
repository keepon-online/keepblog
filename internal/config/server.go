package config

import (
	"time"

	"gitee.com/jieepre/keepblog/config"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	AdminPort    string        `yaml:"adminPort" default:":8000"`
	ReadTimeout  time.Duration `yaml:"readTimeout" default:"30s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" default:"0s"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" default:"2m"`
	Mode         string        `yaml:"mode" default:"release"` // gin模式: debug, release, test
	EnableGzip   bool          `yaml:"enableGzip" default:"true"`
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		AdminPort: ":" + config.Get().Http.Port,
		// 读超时给足图片上传等大体量请求的余量
		ReadTimeout: 30 * time.Second,
		// WriteTimeout 必须为 0：net/http 的写超时是对整个响应的绝对期限，
		// 会把超过时限的流式响应（AI SSE 生成、WebSocket 通知）直接掐断，
		// 客户端表现为 ERR_INCOMPLETE_CHUNKED_ENCODING / ws 连接失败。
		// 长连接的边界由各自机制保证：AI 上游超时（ai.timeout）、客户端
		// 断开时的 ctx 取消，以及下方的 IdleTimeout 兜底连接回收。
		WriteTimeout: 0,
		// 无写超时后，空闲连接由该超时回收，避免连接堆积
		IdleTimeout: 120 * time.Second,
		Mode:        "release",
		EnableGzip:  true,
	}
}
