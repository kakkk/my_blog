package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultFreecacheSize = 10 // MB

type StorageSetting struct {
	Name             string
	KeyFormat        string            `yaml:"key_format"`
	Expire           *int              `yaml:"expire"`
	CacheXSetting    *CacheXSetting    `yaml:"cachex_setting"`
	FreecacheSetting *FreecacheSetting `yaml:"freecache_setting"`
	BigcacheSetting  *BigcacheSetting  `yaml:"bigcache_setting"`
	RedisSetting     *RedisSetting     `yaml:"redis_setting"`
}

func (s *StorageSetting) GetCacheXSetting() *CacheXSetting {
	if s.CacheXSetting == nil {
		return &CacheXSetting{}
	}
	return s.CacheXSetting
}

func (s *StorageSetting) GetExpire() time.Duration {
	if s.Expire == nil {
		return 0
	}
	exp := *s.Expire
	return time.Duration(exp) * time.Minute
}

func (s *StorageSetting) GetFreecacheSetting() *FreecacheSetting {
	if s.FreecacheSetting == nil {
		return &FreecacheSetting{}
	}
	return s.FreecacheSetting
}

func (s *StorageSetting) GetBigcacheSetting() *BigcacheSetting {
	if s.BigcacheSetting == nil {
		return &BigcacheSetting{}
	}
	return s.BigcacheSetting
}

func (s *StorageSetting) GetRedisSetting() *RedisSetting {
	if s.RedisSetting == nil {
		return &RedisSetting{}
	}
	return s.RedisSetting
}

type CacheType string

const (
	CacheTypeRedis     = "redis"
	CacheTypeFreecache = "freecache"
	CacheTypeBigcache  = "bigcache"
)

type CacheXSetting struct {
	AllowDowngrade  *bool       `yaml:"allow_downgrade"`
	DowngradeExpire *int        `yaml:"downgrade_expire"`
	Caches          []CacheType `yaml:"caches"`
}

func (c *CacheXSetting) IsAllowDowngrade() bool {
	if c.AllowDowngrade != nil {
		return *c.AllowDowngrade
	}
	return false
}

func (c *CacheXSetting) GetDowngradeExpire() time.Duration {
	if c.DowngradeExpire == nil {
		return 0
	}
	exp := *c.DowngradeExpire
	return time.Duration(exp) * time.Minute
}

func (c *CacheXSetting) GetCaches() []CacheType {
	return c.Caches
}

type FreecacheSetting struct {
	Size *int `yaml:"size"`
	TTL  *int `yaml:"ttl"`
}

func (f *FreecacheSetting) GetSize() int {
	if f.Size == nil {
		return defaultFreecacheSize * 1024 * 1024
	}
	return *f.Size * 1024 * 1024
}

func (f *FreecacheSetting) GetTTL() time.Duration {
	if f.TTL == nil {
		return 0
	}
	ttl := *f.TTL
	return time.Duration(ttl) * time.Minute
}

type BigcacheSetting struct {
	TTL *int `yaml:"ttl"`
}

func (b *BigcacheSetting) GetTTL() time.Duration {
	if b.TTL == nil {
		return 0
	}
	ttl := *b.TTL
	return time.Duration(ttl) * time.Minute
}

type RedisSetting struct {
	TTL *int `yaml:"ttl"`
}

func (r *RedisSetting) GetTTL() time.Duration {
	if r.TTL == nil {
		return 0
	}
	ttl := *r.TTL
	return time.Duration(ttl) * time.Minute
}

var storageConfig map[string]*StorageSetting

func InitStorageConfig() error {
	storageConfig = make(map[string]*StorageSetting)
	f, err := os.ReadFile("../conf/storage.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, &storageConfig)
	if err != nil {
		return err
	}
	for name, setting := range storageConfig {
		setting.Name = name
	}
	return nil
}

func GetStorageSettingByName(name string) *StorageSetting {
	cfg, ok := storageConfig[name]
	if !ok {
		panic(fmt.Sprintf("config: %v not found", name))
	}
	return cfg
}
