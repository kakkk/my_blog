package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func ServeAdminMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		if strings.HasPrefix(path, "/admin") {
			c.File("../admin/index.html")
			c.Abort()
		} else {
			c.Next(ctx)
		}
	}
}
