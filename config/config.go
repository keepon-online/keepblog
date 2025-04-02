package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
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

// Config file found and successfully parsed
func init() {
	fmt.Println("初始化配置")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./")     // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Config file not found; ignore error if desired
		}
	}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(cfg); err != nil {
			panic(err)
		}
	})
}

func Get() Configs {
	return *cfg
}
