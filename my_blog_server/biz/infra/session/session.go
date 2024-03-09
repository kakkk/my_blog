package session

import (
	"fmt"
	"net/http"

	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/redis"

	"my_blog/biz/infra/config"
)

var sessionStore redis.Store

func MustInit() {
	err := initSessionStore()
	if err != nil {
		panic(err)
	}
}

func initSessionStore() error {
	redisAddr := fmt.Sprintf("%v:%v", config.GetRedisConfig().Host, config.GetRedisConfig().Port)
	s, err := redis.NewStore(20, "tcp", redisAddr, config.GetRedisConfig().Password, []byte("my_blog_session"))
	if err != nil {
		return err
	}
	s.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1天
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	sessionStore = s
	return nil
}

func GetSessionStore() redis.Store {
	return sessionStore
}
