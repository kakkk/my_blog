package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultBlogName        = "kakkk'blog"       // 默认博客名称
	DefaultBlogSubTitle    = "小朱的博客"            // 默认子标题
	DefaultBlogDescription = "这是kakkk（小朱）的技术博客" // 默认博客描述
	DefaultUserName        = "kakkk"            // 默认显示坐着（兜底）
	pageListSize           = 5                  // 列表页大小
)

var appConf *ApplicationConfig

type ApplicationConfig struct {
	App   *App   `yaml:"app"`
	Site  *Site  `yaml:"site"`
	MySQL *MySQL `yaml:"mysql"`
	Redis *Redis `yaml:"redis"`
}

type Site struct {
	SiteDomain string `yaml:"site_domain"`
	CDNDomain  string `yaml:"cdn_domain"`
}

type App struct {
	Name     string `yaml:"name"`
	Release  bool   `yaml:"release"`
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"log_level"`
	LogPath  string `yaml:"log_path"`
}

type MySQL struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func InitApplicationConfig() error {
	appConf = &ApplicationConfig{}
	f, err := os.ReadFile("../conf/application.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, appConf)
	return err
}

func GetAppConfig() *App {
	return appConf.App
}

func GetMySQLConfig() *MySQL {
	return appConf.MySQL
}

func GetRedisConfig() *Redis {
	return appConf.Redis
}

func GetSiteConfig() *Site {
	return appConf.Site
}

func GetPageListSize() int {
	return pageListSize
}

func GetPageListSizeI64() int64 {
	return int64(GetPageListSize())
}

func GetBlogName() string {
	return DefaultBlogName
}

func GetBlogSubTitle() string {
	return DefaultBlogSubTitle
}

func GetDefaultUserName() string {
	return DefaultUserName
}

func GetBlogDescription() string {
	return DefaultBlogDescription
}
