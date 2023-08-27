package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/entity"
	"my_blog/biz/repository/mysql"
)

var userEntityStorage *UserEntityStorage

type UserEntityStorage struct {
	cacheX *cachex.CacheX[int64, *entity.User]
	expire time.Duration
}

func GetUserEntityStorage() *UserEntityStorage {
	return userEntityStorage
}

func initUserEntityStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("user_entity")
	cache, err := NewCacheXBuilderByConfig[int64, *entity.User](ctx, cfg).
		SetGetRealData(userEntityStorageGetRealData).
		SetMGetRealData(userEntityStorageMGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	userEntityStorage = &UserEntityStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func userEntityStorageGetRealData(ctx context.Context, id int64) (*entity.User, error) {
	db := mysql.GetDB(ctx)
	// 获取post
	user, err := mysql.SelectUserByID(db, id)
	if err != nil {
		return parseSqlError(user, err)
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

func (a *UserEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*entity.User {
	return a.cacheX.MGet(ctx, ids, a.expire)
}

func (a *UserEntityStorage) Get(ctx context.Context, id int64) (*entity.User, error) {
	user, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return user, nil
}
