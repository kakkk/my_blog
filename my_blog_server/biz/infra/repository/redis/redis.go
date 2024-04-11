package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"my_blog/biz/infra/config"
)

var (
	client *redis.Client
)

func InitRedis() (err error) {
	cfg := config.GetRedisConfig()
	client = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password, // no password set
		DB:              cfg.DB,       // use default DB
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

func GetClient() *redis.Client {
	return client
}
