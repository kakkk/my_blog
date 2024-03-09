package cachex

import (
	"context"
	"fmt"
	"strings"

	"github.com/kakkk/cachex"
	"github.com/kakkk/cachex/cache"

	"my_blog/biz/infra/config"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/repository/redis"
)

func NewCacheXBuilderByConfig[K comparable, V any](ctx context.Context, cfg *config.StorageSetting) *cachex.Builder[K, V] {
	cxSetting := cfg.GetCacheXSetting()
	builder := cachex.NewBuilder[K, V](ctx).
		SetName(cfg.Name).
		SetLogger(log.NewCacheXLogger()).
		SetGetDataKey(newGetDataKeyFunc[K](cfg.KeyFormat))

	// caches
	for _, cacheType := range cxSetting.GetCaches() {
		switch cacheType {
		case config.CacheTypeFreecache:
			builder.AddCache(newFreecacheByConfig[V](cfg.GetFreecacheSetting()))
		case config.CacheTypeBigcache:
			builder.AddCache(newBigcacheByConfig[V](cfg.GetBigcacheSetting()))
		case config.CacheTypeRedis:
			builder.AddCache(newRedisByConfig[V](cfg.GetRedisSetting()))
		}
	}

	// downgrade
	if cxSetting.IsAllowDowngrade() {
		builder.SetAllowDowngrade(true)
		builder.SetDowngradeCacheExpireTime(cxSetting.GetDowngradeExpire())
	}

	return builder
}

func newFreecacheByConfig[T any](cfg *config.FreecacheSetting) cache.Cache[T] {
	return cache.NewFreeCache[T](cfg.GetSize(), cfg.GetTTL())
}

func newRedisByConfig[T any](cfg *config.RedisSetting) cache.Cache[T] {
	return cache.NewRedisCacheWithClient[T](redis.GetClient(), cfg.GetTTL())
}

func newBigcacheByConfig[T any](cfg *config.BigcacheSetting) cache.Cache[T] {
	return cache.NewBigCache[T](cfg.GetTTL())
}

func newGetDataKeyFunc[K comparable](format string) func(k K) string {
	if !strings.Contains(format, "%v") {
		return func(_ K) string {
			return format
		}
	}
	return func(k K) string {
		return fmt.Sprintf(format, k)
	}
}
