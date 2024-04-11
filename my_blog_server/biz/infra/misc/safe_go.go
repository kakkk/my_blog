package misc

import (
	"context"
	"runtime/debug"

	"my_blog/biz/infra/pkg/log"
)

func SafeGo(ctx context.Context, fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.GetLoggerWithCtx(ctx).Errorf("[panic recover] %v\nstack:%v", err, string(debug.Stack()))
			}
		}()
		fn()
	}()
}
