package repository

import (
	"fmt"

	"my_blog/biz/infra/repository/bleve"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/infra/repository/redis"
)

func MustInit() {
	err := mysql.InitMySQL()
	if err != nil {
		panic(fmt.Sprintf("init mysql fail: %v", err))
	}
	err = redis.InitRedis()
	if err != nil {
		panic(fmt.Sprintf("init redis fail: %v", err))
	}
	err = bleve.Init()
	if err != nil {
		panic(fmt.Sprintf("init bleve fail: %v", err))
	}
}
