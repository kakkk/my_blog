package log

import (
	"context"
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func InitLogger(path string, level string) (err error) {

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("parse level error: %v", err)
	}
	logger = logrus.New()

	// info日志分割
	infoPath := path + "/info.log"
	infoWriter, err := rotatelogs.New(
		infoPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(infoPath),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		return fmt.Errorf("init rotate logs error: %v", err)
	}
	logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel: infoWriter,
		},
		&logrus.JSONFormatter{},
	))

	// 其他日志直接输出
	pathMap := lfshook.PathMap{
		logrus.DebugLevel: path + "/debug.log",
		logrus.ErrorLevel: path + "/error.log",
		logrus.PanicLevel: path + "/panic.log",
		logrus.WarnLevel:  path + "/warn.log",
	}
	logger.Hooks.Add(lfshook.NewHook(pathMap, &logrus.JSONFormatter{}))

	logrus.SetLevel(logLevel)
	logger.Infof("Init logger success ^_^ ")
	return
}

func GetLoggerWithCtx(ctx context.Context) *logrus.Entry {
	// 从context中获取request_id
	requestId, ok := ctx.Value("request_id").(string)
	if !ok {
		requestId = ""
	}
	return logger.WithFields(logrus.Fields{
		"request_id": requestId,
	})
}

func GetLogger() *logrus.Logger {
	return logger
}
