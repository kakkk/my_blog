package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"

	"my_blog/biz/infra/pkg/log"
)

func HertzLoggerMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 记录开始时间
		start := time.Now()
		path := string(c.Request.URI().Path())
		query := string(c.Request.URI().QueryString())
		c.Next(ctx)
		logger := log.GetLoggerWithCtx(ctx)
		// 结束时间
		end := time.Now()
		// 接口耗时
		latency := end.Sub(start)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.WithFields(logrus.Fields{
				"status":     c.Response.StatusCode(),
				"method":     string(c.Request.Method()),
				"path":       path,
				"query":      query,
				"ip":         c.ClientIP(),
				"user-agent": string(c.UserAgent()),
				"latency":    latency,
			}).Infof("success")
		}
	}
}
