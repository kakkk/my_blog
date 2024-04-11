package log

import (
	"context"
)

type CacheXLogger struct{}

func (CacheXLogger) Debug(ctx context.Context, logs string) {
	GetLoggerWithCtx(ctx).Debug(logs)
}

func (CacheXLogger) Info(ctx context.Context, logs string) {
	GetLoggerWithCtx(ctx).Info(logs)
}

func (CacheXLogger) Warn(ctx context.Context, logs string) {
	GetLoggerWithCtx(ctx).Warn(logs)
}

func (CacheXLogger) Error(ctx context.Context, logs string) {
	GetLoggerWithCtx(ctx).Error(logs)
}

func (CacheXLogger) Debugf(ctx context.Context, format string, v ...any) {
	GetLoggerWithCtx(ctx).Debugf(format, v...)
}

func (CacheXLogger) Infof(ctx context.Context, format string, v ...any) {
	GetLoggerWithCtx(ctx).Infof(format, v...)
}

func (CacheXLogger) Warnf(ctx context.Context, format string, v ...any) {
	GetLoggerWithCtx(ctx).Warnf(format, v...)
}

func (CacheXLogger) Errorf(ctx context.Context, format string, v ...any) {
	GetLoggerWithCtx(ctx).Errorf(format, v...)
}

func NewCacheXLogger() CacheXLogger {
	return CacheXLogger{}
}
