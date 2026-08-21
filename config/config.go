package config

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

// 配置以原子快照对外暴露：热重载时整体替换指针，
// 消除旧实现中 Unmarshal 直写共享结构体的并发读写风险。
var (
	cfgPtr atomic.Pointer[Configs]

	loadMu sync.Mutex
	loaded bool
)

type Configs struct {
	Http    *http    `yaml:"http"`
	Minio   *minio   `yaml:"minio"`
	Baidu   *baidu   `yaml:"baidu"`
	Gitalk  *gitalk  `yaml:"gitalk"`
	Redis   *redis   `yaml:"redis"`
	Jwt     *jwt     `yaml:"jwt"`
	Hashids *hashids `yaml:"hashids"`
}

// jwt 签发配置
type jwt struct {
	Secret string `yaml:"secret"`
}

// hashids 编码配置
type hashids struct {
	Salt string `yaml:"salt"`
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

// Load 读取并解析配置，应用启动时显式调用（幂等，重复调用直接返回）。
// 此前版本在包级 init() 中完成这些动作，import 即产生副作用；改为显式
// 加载后，配置文件缺失不再阻塞任何包的导入。
func Load() error {
	loadMu.Lock()
	defer loadMu.Unlock()
	if loaded {
		return nil
	}
	slog.Infof("加载配置")

	// 设置配置文件名和类型
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置环境变量前缀
	viper.SetEnvPrefix("KEEPBLOG")
	viper.AutomaticEnv()

	// 环境变量映射（绑定失败属于配置期错误，启动时无法恢复，故忽略）
	_ = viper.BindEnv("http.port", "KEEPBLOG_HTTP_PORT")
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET", "KEEPBLOG_JWT_SECRET")
	_ = viper.BindEnv("hashids.salt", "KEEPBLOG_HASHIDS_SALT")
	_ = viper.BindEnv("redis.host", "KEEPBLOG_REDIS_HOST")
	_ = viper.BindEnv("redis.port", "KEEPBLOG_REDIS_PORT")
	_ = viper.BindEnv("redis.password", "KEEPBLOG_REDIS_PASSWORD")
	_ = viper.BindEnv("redis.database", "KEEPBLOG_REDIS_DATABASE")
	_ = viper.BindEnv("redis.enable", "KEEPBLOG_REDIS_ENABLE")
	_ = viper.BindEnv("baidu.url", "KEEPBLOG_BASE_URL")

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

	if err := storeSnapshot(); err != nil {
		return err
	}

	// 监听配置文件变化：解析成功后原子替换快照
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件发生变化: %s\n", e.Name)
		if err := storeSnapshot(); err != nil {
			fmt.Printf("重新加载配置失败: %v\n", err)
		}
	})

	loaded = true
	return nil
}

// storeSnapshot 将 viper 当前值解析到全新结构体并原子发布
func storeSnapshot() error {
	fresh := new(Configs)
	if err := viper.Unmarshal(fresh); err != nil {
		return fmt.Errorf("配置解析失败: %w", err)
	}
	cfgPtr.Store(fresh)
	return nil
}

// Get 返回当前配置快照。未调用 Load 时返回零值配置（各指针字段为 nil，
// 调用方需自行判空；应用正常启动流程下不会出现这种情况）。
func Get() Configs {
	if p := cfgPtr.Load(); p != nil {
		return *p
	}
	return Configs{}
}

// setDefaults 设置默认配置值
func setDefaults() {
	// HTTP 默认配置
	viper.SetDefault("http.port", "8000")

	// Redis 默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.enable", false)

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
	viper.SetDefault("minio.bucketName", "keepblog")

	// Gitalk 默认配置
	viper.SetDefault("gitalk.enable", false)
	viper.SetDefault("gitalk.clientID", "")
	viper.SetDefault("gitalk.clientSecret", "")
	viper.SetDefault("gitalk.repo", "")
	viper.SetDefault("gitalk.owner", "")
	viper.SetDefault("gitalk.admin", []string{})

	// JWT / Hashids 默认配置
	viper.SetDefault("jwt.secret", "")
	viper.SetDefault("hashids.salt", "")
}

// ValidateConfig 验证配置
func ValidateConfig() error {
	cfg := Get()

	// JWT 密钥必须显式提供（config.yaml 的 jwt.secret 或环境变量 JWT_SECRET），
	// 不允许回退到代码内置默认值，避免所有部署共享同一密钥
	if cfg.Jwt == nil || cfg.Jwt.Secret == "" {
		return fmt.Errorf("jwt.secret 不能为空：请在 config.yaml 设置 jwt.secret，或设置环境变量 JWT_SECRET（可用 openssl rand -base64 32 生成）")
	}
	if len(cfg.Jwt.Secret) < 32 {
		slog.Warn("jwt.secret 长度小于 32 字符，建议使用更长的随机密钥")
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
