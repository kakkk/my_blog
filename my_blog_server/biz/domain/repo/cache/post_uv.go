package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"my_blog/biz/consts"
)

type PostUVCache struct{}

var postUVCache = new(PostUVCache)

func GetPostUVCache() *PostUVCache {
	return postUVCache
}

func (s *PostUVCache) Get(ctx context.Context, cli *redis.Client, id int64) (int64, error) {
	val, err := cli.Get(ctx, s.getKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, consts.ErrRecordNotFound
		}
		return 0, fmt.Errorf("redis error: %v", err)
	}
	return cast.ToInt64(val), nil
}

func (s *PostUVCache) Incr(ctx context.Context, cli *redis.Client, id int64) error {
	err := cli.Incr(ctx, s.getKey(id)).Err()
	if err != nil {
		return fmt.Errorf("redis error: %v", err)
	}
	return nil
}

func (s *PostUVCache) SetPostUVNX(ctx context.Context, cli *redis.Client, id int64, uv int64) error {
	err := cli.SetNX(ctx, s.getKey(id), uv, 0).Err()
	if err != nil {
		return fmt.Errorf("redis error: %v", err)
	}
	return nil
}

func (s *PostUVCache) getKey(id int64) string {
	return fmt.Sprintf("my_blog:post_uv:%v", id)
}
