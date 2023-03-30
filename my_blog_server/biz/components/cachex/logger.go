package cachex

import (
	"context"
	"fmt"
	"log"
)

type Logger interface {
	Debug(ctx context.Context, logs string)
	Info(ctx context.Context, logs string)
	Warn(ctx context.Context, logs string)
	Error(ctx context.Context, logs string)
	Debugf(ctx context.Context, format string, v ...any)
	Infof(ctx context.Context, format string, v ...any)
	Warnf(ctx context.Context, format string, v ...any)
	Errorf(ctx context.Context, format string, v ...any)
}

var logger Logger

func init() {
	logger = &defaultLogger{}
}

func SetLogger(l Logger) {
	logger = l
}

type defaultLogger struct{}

func (defaultLogger) Debug(ctx context.Context, logs string) {
	log.Printf("[DEBUG] %v\n", logs)
}

func (defaultLogger) Info(ctx context.Context, logs string) {
	log.Printf("[INFO] %v\n", logs)
}

func (defaultLogger) Warn(ctx context.Context, logs string) {
	log.Printf("[WARN] %v\n", logs)
}

func (defaultLogger) Error(ctx context.Context, logs string) {
	log.Printf("[Error] %v\n", logs)
}

func (defaultLogger) Debugf(ctx context.Context, format string, v ...any) {
	defaultLogger{}.Debug(ctx, fmt.Sprintf(format, v...))
}

func (defaultLogger) Infof(ctx context.Context, format string, v ...any) {
	defaultLogger{}.Info(ctx, fmt.Sprintf(format, v...))
}

func (defaultLogger) Warnf(ctx context.Context, format string, v ...any) {
	defaultLogger{}.Warn(ctx, fmt.Sprintf(format, v...))
}

func (defaultLogger) Errorf(ctx context.Context, format string, v ...any) {
	defaultLogger{}.Error(ctx, fmt.Sprintf(format, v...))
}
