package utils

import (
	"context"
	"runtime/debug"

	"my_blog/biz/common/log"
)

func Recover(ctx context.Context, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.GetLoggerWithCtx(ctx).Errorf("[panic recover] %v\nstack:%v", err, string(debug.Stack()))
			fn()
		}
	}()
}
