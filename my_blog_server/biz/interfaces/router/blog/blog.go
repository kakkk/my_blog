// Code generated by hertz generator. DO NOT EDIT.

package blog

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	blog "my_blog/biz/interfaces/handler/blog"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	root.GET("/", append(_indexpageMw(), blog.IndexPage)...)
	root.GET("/archives", append(_archivespageMw(), blog.ArchivesPage)...)
	_archives := root.Group("/archives", _archivesMw()...)
	_archives.GET("/:post_id", append(_postpageMw(), blog.PostPage)...)
	root.GET("/categories", append(_categoriespageMw(), blog.CategoriesPage)...)
	root.GET("/search", append(_searchpageMw(), blog.SearchPage)...)
	root.GET("/tags", append(_tagspageMw(), blog.TagsPage)...)
	{
		_api := root.Group("/api", _apiMw()...)
		_api.GET("/captcha", append(_getcaptchaapiMw(), blog.GetCaptchaAPI)...)
		_api.GET("/search", append(_searchapiMw(), blog.SearchAPI)...)
		{
			_admin := _api.Group("/admin", _adminMw()...)
			_admin.POST("/login", append(_loginapiMw(), blog.LoginAPI)...)
			_admin.POST("/category", append(_createcategoryapiMw(), blog.CreateCategoryAPI)...)
			_category := _admin.Group("/category", _categoryMw()...)
			_category.PUT("/:category_id", append(_updatecategoryapiMw(), blog.UpdateCategoryAPI)...)
			_category.DELETE("/:category_id", append(_deletecategoryapiMw(), blog.DeleteCategoryAPI)...)
			_category.GET("/list", append(_getcategorylistapiMw(), blog.GetCategoryListAPI)...)
			_category.PUT("/order", append(_updatecategoryorderapiMw(), blog.UpdateCategoryOrderAPI)...)
			{
				_comment := _admin.Group("/comment", _commentMw()...)
				_comment.PUT("/:comment_id", append(_deletecommentadminapiMw(), blog.DeleteCommentAdminAPI)...)
				_comment.GET("/list", append(_getcommentlistadminapiMw(), blog.GetCommentListAdminAPI)...)
				{
					_comment_id := _comment.Group("/:comment_id", _comment_idMw()...)
					_comment_id.POST("/reply", append(_replycommentadminapiMw(), blog.ReplyCommentAdminAPI)...)
					_comment_id.PUT("/status", append(_updatecommentstatusadminapiMw(), blog.UpdateCommentStatusAdminAPI)...)
				}
			}
			_admin.POST("/post", append(_createpostapiMw(), blog.CreatePostAPI)...)
			_post := _admin.Group("/post", _postMw()...)
			_post.PUT("/:post_id", append(_updatepostapiMw(), blog.UpdatePostAPI)...)
			_post.DELETE("/:post_id", append(_deletepostapiMw(), blog.DeletePostAPI)...)
			_post.POST("/list", append(_getpostlistapiMw(), blog.GetPostListAPI)...)
			_post.GET("/:post_id", append(_getpostapiMw(), blog.GetPostAPI)...)
			_post_id := _post.Group("/:post_id", _post_idMw()...)
			_post_id.PUT("/status", append(_updatepoststatusapiMw(), blog.UpdatePostStatusAPI)...)
			_admin.POST("/tag", append(_createtagapiMw(), blog.CreateTagAPI)...)
			_tag := _admin.Group("/tag", _tagMw()...)
			_tag.PUT("/:tag_id", append(_updatetagapiMw(), blog.UpdateTagAPI)...)
			_tag.DELETE("/:tag_id", append(_deletetagapiMw(), blog.DeleteTagAPI)...)
			_tag.GET("/list", append(_gettaglistapiMw(), blog.GetTagListAPI)...)
			{
				_user := _admin.Group("/user", _userMw()...)
				_user.GET("/info", append(_getuserinfoapiMw(), blog.GetUserInfoAPI)...)
			}
		}
		{
			_comment0 := _api.Group("/comment", _comment0Mw()...)
			_comment0.POST("/article", append(_commentarticleapiMw(), blog.CommentArticleAPI)...)
			_comment0.GET("/list", append(_getcommentlistapiMw(), blog.GetCommentListAPI)...)
			_comment0.POST("/reply", append(_replycommentapiMw(), blog.ReplyCommentAPI)...)
		}
	}
	{
		_category0 := root.Group("/category", _category0Mw()...)
		_category0.GET("/:name", append(_categorypostpageMw(), blog.CategoryPostPage)...)
		_name := _category0.Group("/:name", _nameMw()...)
		_name.GET("/:page", append(_categorypostbypaginationpageMw(), blog.CategoryPostByPaginationPage)...)
	}
	{
		_page := root.Group("/page", _pageMw()...)
		_page.GET("/:page", append(_indexbypaginationpageMw(), blog.IndexByPaginationPage)...)
	}
	{
		_tag0 := root.Group("/tag", _tag0Mw()...)
		_tag0.GET("/:name", append(_tagpostpageMw(), blog.TagPostPage)...)
		_name0 := _tag0.Group("/:name", _name0Mw()...)
		_name0.GET("/:page", append(_tagpostbypaginationpageMw(), blog.TagPostByPaginationPage)...)
	}
}
