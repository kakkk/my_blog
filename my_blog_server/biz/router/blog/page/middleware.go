// Code generated by hertz generator.

package Page

import (
	"my_blog/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	// your code...
	mw := middleware.GetRootMW()
	mw = append(mw, middleware.VisitorSessionMW())
	return mw
}

func _indexpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _categoriespageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tagspageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _categoryMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tagMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _slugMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _nameMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _categorypostbypaginationpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _pageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _indexbypaginationpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tagpostbypaginationpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _searchpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _archivesMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _postpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _name0Mw() []app.HandlerFunc {
	// your code...
	return nil
}
