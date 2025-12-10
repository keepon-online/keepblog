package config

import (
	"time"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	AdminPort    string        `yaml:"adminPort" default:":8000"`
	ConsolePort  string        `yaml:"consolePort" default:":8890"`
	WebPort      string        `yaml:"webPort" default:":8589"`
	ReadTimeout  time.Duration `yaml:"readTimeout" default:"5s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" default:"10s"`
	Mode         string        `yaml:"mode" default:"release"` // gin模式: debug, release, test
	EnableGzip   bool          `yaml:"enableGzip" default:"true"`
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		AdminPort:    ":8001",
		ConsolePort:  ":8890",
		WebPort:      ":8589",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Mode:         "release",
		EnableGzip:   true,
	}
}
