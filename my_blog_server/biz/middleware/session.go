package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"

	"my_blog/biz/infra/session"
)

func SessionMW() app.HandlerFunc {
	return sessions.Many([]string{"session", "admin_session"}, session.GetSessionStore())
}
