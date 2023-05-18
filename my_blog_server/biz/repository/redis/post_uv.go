package redis

import (
	"errors"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/cast"

	"my_blog/biz/common/consts"
)

func GetPostUV(cli *redis.Client, key string) (int64, error) {
	val, err := cli.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, consts.ErrRecordNotFound
		}
		return 0, err
	}
	return cast.ToInt64(val), nil
}

func IncrPostUV(cli *redis.Client, key string) error {
	return cli.Incr(key).Err()
}

func SetPostUVNX(cli *redis.Client, key string, uv int64) error {
	return cli.SetNX(key, uv, 0).Err()
}
