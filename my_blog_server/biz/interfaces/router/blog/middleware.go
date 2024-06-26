// Code generated by hertz generator.

package blog

import (
	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/interfaces/middleware"
)

func rootMw() []app.HandlerFunc {
	return middleware.GetRootMW()
}

func _indexpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _archivesMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _archivespageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _postpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _categoriespageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _searchpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tagspageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _apiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getcaptchaapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _searchapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _adminMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.AdminSessionMW(),
	}
}

func _loginapiMw() []app.HandlerFunc { return nil }

func _categoryMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _createcategoryapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updatecategoryapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletecategoryapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getcategorylistapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updatecategoryorderapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commentMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletecommentadminapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getcommentlistadminapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _comment_idMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _replycommentadminapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updatecommentstatusadminapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _postMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _createpostapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updatepostapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletepostapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getpostlistapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _post_idMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getpostapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updatepoststatusapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tagMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _createtagapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updatetagapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletetagapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _gettaglistapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getuserinfoapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _comment0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commentarticleapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getcommentlistapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _replycommentapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _category0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _nameMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _categorypostpageMw() []app.HandlerFunc {
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

func _tag0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _name0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tagpostpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tagpostbypaginationpageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _createpageapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getpageapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _updatepageapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deletepageapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getpagelistapiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _page0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _pagesMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _pagepageMw() []app.HandlerFunc {
	// your code...
	return nil
}
