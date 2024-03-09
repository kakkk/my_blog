package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/infra/config"
	cachex2 "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/model"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/mysql"
)

var userEntityStorage *UserEntityStorage

type UserEntityStorage struct {
	cacheX *cachex.CacheX[int64, *model.User]
	expire time.Duration
}

func GetUserEntityStorage() *UserEntityStorage {
	return userEntityStorage
}

func initUserEntityStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("user_entity")
	cache, err := cachex2.NewCacheXBuilderByConfig[int64, *model.User](ctx, cfg).
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

func userEntityStorageGetRealData(ctx context.Context, id int64) (*model.User, error) {
	db := mysql2.GetDB(ctx)
	// 获取post
	user, err := mysql.SelectUserByID(db, id)
	if err != nil {
		return parseSqlError(user, err)
	}
	return user, nil
}

func userEntityStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*model.User, error) {
	db := mysql2.GetDB(ctx)
	users, err := mysql.MSelectUserByIDs(db, ids)
	if err != nil {
		return parseSqlError(users, err)
	}
	return users, nil
}

func (a *UserEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*model.User {
	return a.cacheX.MGet(ctx, ids, a.expire)
}

func (a *UserEntityStorage) Get(ctx context.Context, id int64) (*model.User, error) {
	user, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return user, nil
}
