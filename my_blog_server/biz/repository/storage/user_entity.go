package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
	"my_blog/biz/entity"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var userEntityStorage *UserEntityStorage

type UserEntityStorage struct {
	cacheX *cachex.CacheX[*entity.User, int64]
}

func GetUserEntityStorage() *UserEntityStorage {
	return userEntityStorage
}

func initUserEntityStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache(ctx, 1000, time.Minute)
	cache := cachex.NewSerializableCacheX[*entity.User, int64]("user_entity", false, true).
		SetGetCacheKey(userEntityStorageGetKey).
		SetGetRealData(userEntityStorageGetRealData).
		SetMGetRealData(userEntityStorageMGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	userEntityStorage = &UserEntityStorage{
		cacheX: cache,
	}
	return nil
}

func userEntityStorageGetRealData(ctx context.Context, id int64) (*entity.User, error) {
	db := mysql.GetDB(ctx)
	// 获取post
	user, err := mysql.SelectUserByID(db, id)
	if err != nil {
		if err == consts.ErrRecordNotFound {
			return nil, cachex.ErrNotFound
		}
		return nil, fmt.Errorf("sql error: %w", err)
	}
	return user, nil
}

func userEntityStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*entity.User, error) {
	db := mysql.GetDB(ctx)
	users, err := mysql.MSelectUserByIDs(db, ids)
	if err != nil {
		return parseSqlError(users, err)
	}
	return users, nil
}

func userEntityStorageGetKey(id int64) string {
	return fmt.Sprintf("my_blog_user_entity_%v", id)
}

func (a *UserEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*entity.User {
	return a.cacheX.MGet(ctx, ids)
}

func (a *UserEntityStorage) Get(ctx context.Context, id int64) (*entity.User, error) {
	article, err := a.cacheX.Get(ctx, id)
	if err != nil {
		if errors.Is(err, cachex.ErrNotFound) {
			return nil, consts.ErrRecordNotFound
		}
		return nil, fmt.Errorf("cacheX error:[%w]", err)
	}

	return article, nil
}