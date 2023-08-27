package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"my_blog/biz/common/consts"
)

func GetPostUV(ctx context.Context, cli *redis.Client, key string) (int64, error) {
	val, err := cli.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, consts.ErrRecordNotFound
		}
		return 0, err
	}
	return cast.ToInt64(val), nil
}

func IncrPostUV(ctx context.Context, cli *redis.Client, key string) error {
	return cli.Incr(ctx, key).Err()
}

func SetPostUVNX(ctx context.Context, cli *redis.Client, key string, uv int64) error {
	return cli.SetNX(ctx, key, uv, 0).Err()
}
