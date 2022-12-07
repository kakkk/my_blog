package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var conf *Config

type Config struct {
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

func InitConfig() error {
	conf = &Config{}
	f, err := ioutil.ReadFile("../conf/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, conf)
	return err
}

func GetAppConfig() *App {
	return conf.App
}

func GetMySQLConfig() *MySQL {
	return conf.MySQL
}

func GetRedisConfig() *Redis {
	return conf.Redis
}

func GetSiteConfig() *Site {
	return conf.Site
}
