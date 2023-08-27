package redis

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/common/config"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

func InitRedis() (err error) {
	conf := config.GetRedisConfig()
	client = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password:        conf.Password, // no password set
		DB:              conf.DB,       // use default DB
		ReadTimeout:     500 * time.Millisecond,
		WriteTimeout:    500 * time.Millisecond,
		ConnMaxIdleTime: 60 * time.Second,
		PoolSize:        64,
		MinIdleConns:    16,
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
