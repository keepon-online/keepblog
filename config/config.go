package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var cfg = new(Configs)

type Configs struct {
	Http   *http   `yaml:"http"`
	Mysql  *mysql  `yaml:"mysql"`
	Minio  *minio  `yaml:"minio"`
	System *system `yaml:"system"`
	Baidu  *baidu  `yaml:"system"`
	Gitalk *gitalk `yaml:"gitalk"`
	Redis  *redis  `yaml:"redis"`
}

type system struct {
	BaseUrl string `yaml:"baseUrl"`
}

// 百度收录
type baidu struct {
	Push  bool   `yaml:"push"`
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

// 服务绑定地址
type http struct {
	Port string `yaml:"port"`
}

// mysql 配置
type mysql struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
}

// minio配置
type minio struct {
	ServerUrl       string `yaml:"serverUrl"`
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `yaml:"useSSL"`
	BucketName      string `yaml:"bucketName"`
}
type gitalk struct {
	Enable       bool     `yaml:"enable"`
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	Repo         string   `yaml:"repo"`
	Owner        string   `yaml:"owner"`
	Admin        []string `yaml:"admin"`
}

// redis配置
type redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
	Enable   bool   `yaml:"enable"`
}

// Config file found and successfully parsed
func init() {
	slog.Infof("初始化配置")

	// 设置配置文件名和类型
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置环境变量前缀
	viper.SetEnvPrefix("GOSITE")
	viper.AutomaticEnv()

	// 环境变量映射
	viper.BindEnv("http.port", "GOSITE_HTTP_PORT")
	viper.BindEnv("mysql.host", "GOSITE_MYSQL_HOST")
	viper.BindEnv("mysql.port", "GOSITE_MYSQL_PORT")
	viper.BindEnv("mysql.username", "GOSITE_MYSQL_USERNAME")
	viper.BindEnv("mysql.password", "GOSITE_MYSQL_PASSWORD")
	viper.BindEnv("mysql.database", "GOSITE_MYSQL_DATABASE")
	viper.BindEnv("redis.host", "GOSITE_REDIS_HOST")
	viper.BindEnv("redis.port", "GOSITE_REDIS_PORT")
	viper.BindEnv("redis.password", "GOSITE_REDIS_PASSWORD")
	viper.BindEnv("redis.database", "GOSITE_REDIS_DATABASE")
	viper.BindEnv("redis.enable", "GOSITE_REDIS_ENABLE")
	viper.BindEnv("system.baseUrl", "GOSITE_BASE_URL")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("配置文件未找到，使用默认配置和环境变量")
		} else {
			fmt.Printf("读取配置文件出错: %v\n", err)
		}
	}

	// 解析配置
	if err := viper.Unmarshal(cfg); err != nil {
		panic(fmt.Sprintf("配置解析失败: %v", err))
	}

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件发生变化: %s\n", e.Name)
		if err := viper.Unmarshal(cfg); err != nil {
			fmt.Printf("重新加载配置失败: %v\n", err)
		}
	})
}

func Get() Configs {
	return *cfg
}

// setDefaults 设置默认配置值
func setDefaults() {
	// HTTP 默认配置
	viper.SetDefault("http.port", "8000")

	// MySQL 默认配置
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.username", "root")
	viper.SetDefault("mysql.password", "")
	viper.SetDefault("mysql.database", "go_site")

	// Redis 默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.enable", false)

	// 系统默认配置
	viper.SetDefault("system.baseUrl", "http://localhost:8000")

	// Baidu 默认配置
	viper.SetDefault("baidu.push", false)
	viper.SetDefault("baidu.url", "")
	viper.SetDefault("baidu.token", "")

	// Minio 默认配置
	viper.SetDefault("minio.serverUrl", "")
	viper.SetDefault("minio.endpoint", "")
	viper.SetDefault("minio.accessKeyID", "")
	viper.SetDefault("minio.secretAccessKey", "")
	viper.SetDefault("minio.useSSL", false)
	viper.SetDefault("minio.bucketName", "go-site")

	// Gitalk 默认配置
	viper.SetDefault("gitalk.enable", false)
	viper.SetDefault("gitalk.clientID", "")
	viper.SetDefault("gitalk.clientSecret", "")
	viper.SetDefault("gitalk.repo", "")
	viper.SetDefault("gitalk.owner", "")
	viper.SetDefault("gitalk.admin", []string{})
}

// ValidateConfig 验证配置
func ValidateConfig() error {
	cfg := Get()

	// 验证基本配置
	if cfg.System == nil || cfg.System.BaseUrl == "" {
		return fmt.Errorf("system.baseUrl 不能为空")
	}

	// 验证Redis配置（如果启用）
	if cfg.Redis != nil && cfg.Redis.Enable {
		if cfg.Redis.Host == "" {
			return fmt.Errorf("redis.host 不能为空")
		}
		if cfg.Redis.Port <= 0 || cfg.Redis.Port > 65535 {
			return fmt.Errorf("redis.port 必须在 1-65535 范围内")
		}
	}

	// 验证Minio配置（如果配置了）
	if cfg.Minio != nil && cfg.Minio.Endpoint != "" {
		if cfg.Minio.AccessKeyID == "" || cfg.Minio.SecretAccessKey == "" {
			return fmt.Errorf("minio accessKeyID 和 secretAccessKey 不能为空")
		}
	}

	return nil
}
