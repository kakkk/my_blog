// Code generated by hertz generator.

package main

import (
	"html/template"

	"my_blog/biz/handler"
	"my_blog/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {
	r.SetFuncMap(map[string]any{
		"unescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	})
	r.LoadHTMLGlob("../templates/*")
	r.Static("/assets", "../")
	r.NoRoute(append(middleware.GetRootMW(), handler.NotFoundHandler)...)
	r.GET("/ping", handler.Ping)

}