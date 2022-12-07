package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
)

func HertzRecoverMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				logger := log.GetLoggerWithCtx(ctx)
				logger.WithFields(logrus.Fields{
					"method": string(c.GetRequest().Method()),
					"path":   string(c.GetRequest().Path()),
					"header": string(c.GetRequest().Header.Header()),
					"body":   string(c.GetRequest().Body()),
					"error":  err,
					"stack":  string(debug.Stack()),
				}).Error("[Recovery from panic]")
				c.HTML(http.StatusInternalServerError, "index.tmpl", resp.GetInternalErrorPageResp())
				c.Abort()
			}
		}()
		c.Next(ctx)
	}
}
