package config

import (
	"time"

	"gitee.com/jieepre/keepblog/config"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	AdminPort    string        `yaml:"adminPort" default:":8000"`
	ReadTimeout  time.Duration `yaml:"readTimeout" default:"5s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" default:"10s"`
	Mode         string        `yaml:"mode" default:"release"` // gin模式: debug, release, test
	EnableGzip   bool          `yaml:"enableGzip" default:"true"`
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		AdminPort:    ":" + config.Get().Http.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Mode:         "release",
		EnableGzip:   true,
	}
}
