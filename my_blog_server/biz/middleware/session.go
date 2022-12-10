package middleware

import (
	"context"
	"fmt"
	"net/http"

	"my_blog/biz/common/config"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/redis"
)

var store redis.Store

func InitSession() error {
	redisAddr := fmt.Sprintf("%v:%v", config.GetRedisConfig().Host, config.GetRedisConfig().Port)
	s, err := redis.NewStore(10, "tcp", redisAddr, config.GetRedisConfig().Password, []byte("my_blog_session"))
	if err != nil {
		return err
	}
	s.Options(sessions.Options{
		Path:   "/",
		MaxAge: 604800, // 7天
	})
	store = s
	return nil
}

func SessionMW() app.HandlerFunc {
	return sessions.New("session", store)
}

func SessionAuthMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		logger := log.GetLoggerWithCtx(ctx)
		session := sessions.Default(c)
		userIDFromSession := session.Get("user_id")
		if userIDFromSession == nil {
			c.JSON(http.StatusUnauthorized, resp.NewBaseResponse(common.RespCode_Unauthorized, ""))
			c.Abort()
			return
		}
		userID, ok := userIDFromSession.(int64)
		if !ok {
			logger.Warnf("convert user_id fail")
			c.JSON(http.StatusUnauthorized, resp.NewBaseResponse(common.RespCode_Unauthorized, ""))
			c.Abort()
			return
		}
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, resp.NewBaseResponse(common.RespCode_Unauthorized, ""))
			c.Abort()
			return
		}
		logger.Infof("session user_id:[%v]", userID)
		ctx = context.WithValue(ctx, "user_id", userID)
	}
}
