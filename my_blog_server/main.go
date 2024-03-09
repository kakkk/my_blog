// Code generated by hertz generator.

package main

import (
	"fmt"

	"my_blog/biz/infra/config"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/repository"
	"my_blog/biz/infra/session"
	"my_blog/biz/repository/index"
	"my_blog/biz/repository/storage"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
)

func main() {
	// 先初始化配置
	config.MustInit()

	// 日志
	log.MustInit()

	// 存储
	repository.MustInit()

	// session
	session.MustInit()

	// storage
	if err := storage.InitStorage(); err != nil {
		panic(err)
	}

	// index
	if err := index.InitArticleIndex(); err != nil {
		panic(err)
	}

	// init service
	if err := service.InitService(); err != nil {
		panic(err)
	}

	// hertz
	h := initHertz()
	register(h)
	h.Spin()
}

func initHertz() *server.Hertz {
	cfg := config.GetAppConfig()
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("127.0.0.1:%v", cfg.Port)),
	)
	hlog.SetLogger(hertzlogrus.NewLogger(
		hertzlogrus.WithLogger(log.GetLogger()),
	))
	return h
}
