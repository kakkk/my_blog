// Code generated by hertz generator.

package Api

import (
	"my_blog/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return middleware.GetRootMW()
}

func _apiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _login_piMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _adminMw() []app.HandlerFunc {
	// your code...
	return nil
}
