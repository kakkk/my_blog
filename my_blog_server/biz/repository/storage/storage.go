package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/kakkk/cachex"
	"github.com/kakkk/cachex/cache"

	"my_blog/biz/common/config"
	"my_blog/biz/common/log"
	"my_blog/biz/repository/redis"
)

func InitStorage() error {
	ctx := context.Background()
	//cachex.SetLogger(log.NewCacheXLogger())
	err := initArticleEntityStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostOrderListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostPrevNextStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostMetaStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initUserEntityStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategoryEntityStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostTagListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostCategoryListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategoryListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategoryPostListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategorySlugIDStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initTagListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initTagNameIDStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initTagPostListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	initArchivesStorage()
	return nil
}

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
	return cache.NewRedisCacheWithClient[T](redis.GetRedisClient(), cfg.GetTTL())
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
