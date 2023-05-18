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

	logger.SetLevel(logLevel)
	logger.Infof("Init logger success ^_^ ")
	return
}

func GetLoggerWithCtx(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}
	// 从context中获取request_id
	requestID, ok := ctx.Value("request_id").(string)
	if ok && requestID != "" {
		fields["request_id"] = requestID
	}

	sessionID, ok := ctx.Value("session_id").(string)
	if ok && sessionID != "" {
		fields["session_id"] = sessionID
	}

	userID, ok := ctx.Value("user_id").(int64)
	if ok && userID != 0 {
		fields["user_id"] = userID
	}
	return logger.WithFields(fields)
}

func GetLogger() *logrus.Logger {
	return logger
}
