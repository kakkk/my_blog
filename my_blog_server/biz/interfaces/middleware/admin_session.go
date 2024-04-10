package middleware

import (
	"context"
	"net/http"

	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
)

func AdminSessionMW(auth bool) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		logger := log.GetLoggerWithCtx(ctx)
		session := sessions.DefaultMany(c, "admin_session")
		ctx = context.WithValue(ctx, "admin_session", session)
		ctx = context.WithValue(ctx, "session_id", session.ID())
		if !auth {
			c.Next(ctx)
		}
		userID, ok := session.Get("user_id").(int64)
		if !ok {
			logger.Warnf("get user_id from session fail")
			c.JSON(http.StatusUnauthorized, resp.NewBaseResponse(common.RespCode_Unauthorized, ""))
			c.Abort()
			return
		}
		if userID == 0 {
			logger.Warnf("session user_id is 0")
			c.JSON(http.StatusUnauthorized, resp.NewBaseResponse(common.RespCode_Unauthorized, ""))
			c.Abort()
			return
		}
		logger.Infof("session user_id:[%v]", userID)
		ctx = context.WithValue(ctx, "user_id", userID)
		c.Next(ctx)
		err := session.Save()
		if err != nil {
			log.GetLoggerWithCtx(ctx).Warnf("save session error:[%v]", err)
		}
	}
}
