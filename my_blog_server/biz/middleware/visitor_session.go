package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"

	"my_blog/biz/infra/pkg/log"
)

func VisitorSessionMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.DefaultMany(c, "session")
		ctx = context.WithValue(ctx, "session", session)
		ctx = context.WithValue(ctx, "session_id", session.ID())
		c.Next(ctx)
		err := session.Save()
		if err != nil {
			log.GetLoggerWithCtx(ctx).Errorf("save session error")
		}
	}
}
