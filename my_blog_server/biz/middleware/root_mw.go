package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
)

func GetRootMW() []app.HandlerFunc {
	return []app.HandlerFunc{
		RequestIdMW(),
		HertzLoggerMW(),
		CorsMW(),
		SessionMW(),
	}
}

func GetNoRouteMW() []app.HandlerFunc {
	return []app.HandlerFunc{
		RequestIdMW(),
		HertzLoggerMW(),
		CorsMW(),
		SessionMW(),
		ServeAdminMW(),
		NotFoundMW(),
	}
}
