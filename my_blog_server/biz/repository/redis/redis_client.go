package redis

import (
	"context"
	"fmt"

	"my_blog/biz/common/config"

	"github.com/go-redis/redis/v7"
)

var (
	client *redis.Client
)

func InitRedis() (err error) {
	conf := config.GetRedisConfig()
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
	})
	_, err = client.Ping().Result()
	if err != nil {
		return
	}
	return err
}

func GetRedisClient(ctx context.Context) *redis.Client {
	return client.WithContext(ctx)
}
