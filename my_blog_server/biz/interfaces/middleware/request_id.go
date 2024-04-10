package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

func RequestIdMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestId := c.Request.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = strings.Replace(uuid.New().String(), "-", "", -1)
			c.Request.Header.Set("X-Request-ID", requestId)
		}
		c.Header("X-Request-ID", requestId)
		c.Next(context.WithValue(ctx, "request_id", requestId))
	}
}
