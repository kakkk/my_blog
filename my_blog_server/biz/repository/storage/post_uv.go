package storage

import (
	"context"
	"errors"
	"fmt"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/log"
	"my_blog/biz/repository/redis"
)

type PostUVStorage struct{}

var postUVStorage = new(PostUVStorage)

func GetPostUVStorage() *PostUVStorage {
	return postUVStorage
}

func (s *PostUVStorage) Get(ctx context.Context, id int64) (int64, error) {
	uv, err := redis.GetPostUV(ctx, redis.GetRedisClient(), getPostUVKey(id))
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return s.updateRedis(ctx, id), nil
		}
		return 0, fmt.Errorf("redis error:[%v]", err)
	}
	return uv, nil
}

func (s *PostUVStorage) Incr(ctx context.Context, id int64) error {
	return redis.IncrPostUV(ctx, redis.GetRedisClient(), getPostUVKey(id))
}

func (s *PostUVStorage) updateRedis(ctx context.Context, id int64) int64 {
	// 回源直接从cache拿即可，UV数据实时性要求不高
	logger := log.GetLoggerWithCtx(ctx)
	postEntity, err := GetArticleEntityStorage().Get(ctx, id)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			logger.Warnf("post_id not found, id:[%v]", id)
			return 0
		}
		logger.Errorf("get post entity error:[%v]", err)
		return 0
	}
	err = redis.SetPostUVNX(ctx, redis.GetRedisClient(), getPostUVKey(id), postEntity.UV)
	if err != nil {
		logger.Errorf("set post uv nx error:[%v]", err)
	}
	return postEntity.UV
}

func getPostUVKey(id int64) string {
	return fmt.Sprintf("my_blog:post_uv:%v", id)
}
