package logger

import (
	"context"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func InitLogger(path string, level string) (err error) {
	// 可使用file-rotatelogs进行日志分割
	// 参考：https://github.com/rifflock/lfshook
	pathMap := lfshook.PathMap{
		logrus.DebugLevel: path + "/debug.log",
		logrus.InfoLevel:  path + "/info.log",
		logrus.ErrorLevel: path + "/error.log",
		logrus.PanicLevel: path + "/panic.log",
		logrus.WarnLevel:  path + "/warn.log",
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logger = logrus.New()
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
