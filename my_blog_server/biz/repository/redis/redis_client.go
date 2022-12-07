package redis

import (
	"context"
	"fmt"

	"my_blog/biz/common/config"

	"github.com/go-redis/redis/v8"
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
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return
	}
	return err
}

func GetRedisClient() *redis.Client {
	return client
}
