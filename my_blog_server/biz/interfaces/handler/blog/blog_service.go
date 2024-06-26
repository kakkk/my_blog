// Code generated by hertz generator.

package blog

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	_ "my_blog/biz/hertz_gen/blog/api"
	_ "my_blog/biz/hertz_gen/blog/page"
	"my_blog/biz/interfaces/facade"
)

// LoginAPI .
// @router /api/login [POST]
func LoginAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.LoginAPI(ctx, c))
}

// GetUserInfoAPI .
// @router /api/admin/user/info [GET]
func GetUserInfoAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetUserInfoAPI(ctx, c))
}

// CreateTagAPI .
// @router /api/admin/tag [POST]
func CreateTagAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.CreateTagAPI(ctx, c))
}

// UpdateTagAPI .
// @router /api/admin/tag/:tag_id [PUT]
func UpdateTagAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.UpdateTagAPI(ctx, c))
}

// DeleteTagAPI .
// @router /api/admin/tag/:tag_id [DELETE]
func DeleteTagAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.DeleteTagAPI(ctx, c))
}

// GetTagListAPI .
// @router /api/admin/tag/list [GET]
func GetTagListAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetTagListAPI(ctx, c))
}

// CreateCategoryAPI .
// @router /api/admin/category [POST]
func CreateCategoryAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.CreateCategoryAPI(ctx, c))
}

// UpdateCategoryAPI .
// @router /api/admin/category/:category_id [PUT]
func UpdateCategoryAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.UpdateCategoryAPI(ctx, c))
}

// DeleteCategoryAPI .
// @router /api/admin/category/:category_id [DELETE]
func DeleteCategoryAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.DeleteCategoryAPI(ctx, c))
}

// UpdateCategoryOrderAPI .
// @router /api/admin/category/order [PUT]
func UpdateCategoryOrderAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.UpdateCategoryOrderAPI(ctx, c))
}

// GetCategoryListAPI .
// @router /api/admin/category/list [GET]
func GetCategoryListAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetCategoryListAPI(ctx, c))
}

// CreatePostAPI .
// @router /api/admin/post [POST]
func CreatePostAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.CreatePostAPI(ctx, c))
}

// GetPostAPI .
// @router /api/admin/post [GET]
func GetPostAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetPostAPI(ctx, c))
}

// UpdatePostAPI .
// @router /api/admin/post/:post_id [PUT]
func UpdatePostAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.UpdatePostAPI(ctx, c))
}

// UpdatePostStatusAPI .
// @router /api/admin/post/:post_id/status [PUT]
func UpdatePostStatusAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.UpdatePostStatusAPI(ctx, c))
}

// GetPostListAPI .
// @router /api/admin/post/list [POST]
func GetPostListAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetPostListAPI(ctx, c))
}

// DeletePostAPI .
// @router /api/admin/post/:post_id [DELETE]
func DeletePostAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.DeletePostAPI(ctx, c))
}

// SearchAPI .
// @router /api/search [GET]
func SearchAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.SearchAPI(ctx, c))
}

// GetCommentListAdminAPI .
// @router /api/admin/comment/list [GET]
func GetCommentListAdminAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetCommentListAdminAPI(ctx, c))
}

// ReplyCommentAdminAPI .
// @router /api/admin/comment/:comment_id/reply [POST]
func ReplyCommentAdminAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.ReplyCommentAdminAPI(ctx, c))
}

// UpdateCommentStatusAdminAPI .
// @router /api/admin/comment/:comment_id/status [PUT]
func UpdateCommentStatusAdminAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.UpdateCommentStatusAdminAPI(ctx, c))
}

// DeleteCommentAdminAPI .
// @router /api/admin/comment/:comment_id [PUT]
func DeleteCommentAdminAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.DeleteCommentAdminAPI(ctx, c))
}

// GetCommentListAPI .
// @router /api/comment/list [GET]
func GetCommentListAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetCommentListAPI(ctx, c))
}

// CommentArticleAPI .
// @router /api/comment/article [GET]
func CommentArticleAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.CommentArticleAPI(ctx, c))
}

// ReplyCommentAPI .
// @router /api/comment/reply [GET]
func ReplyCommentAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.ReplyCommentAPI(ctx, c))
}

// GetCaptchaAPI .
// @router /api/captcha [GET]
func GetCaptchaAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetCaptchaAPI(ctx, c))
}

// IndexPage .
// @router / [GET]
func IndexPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.IndexPage(ctx, c))
}

// IndexByPaginationPage .
// @router /page/:page [GET]
func IndexByPaginationPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.IndexByPaginationPage(ctx, c))
}

// CategoryPostPage .
// @router /category/:slug [GET]
func CategoryPostPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.CategoryPostPage(ctx, c))
}

// CategoryPostByPaginationPage .
// @router /category/:slug/:page [GET]
func CategoryPostByPaginationPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.CategoryPostByPaginationPage(ctx, c))
}

// TagPostPage .
// @router /tag/:name [GET]
func TagPostPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.TagPostPage(ctx, c))
}

// TagPostByPaginationPage .
// @router /tag/:name/:page [GET]
func TagPostByPaginationPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.TagPostByPaginationPage(ctx, c))
}

// ArchivesPage .
// @router /archives [GET]
func ArchivesPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.ArchivesPage(ctx, c))
}

// TagsPage .
// @router /tags [GET]
func TagsPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.TagsPage(ctx, c))
}

// CategoriesPage .
// @router /categories [GET]
func CategoriesPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.CategoriesPage(ctx, c))
}

// SearchPage .
// @router /search [GET]
func SearchPage(ctx context.Context, c *app.RequestContext) {
	resp, code := facade.SearchPage(ctx, c)
	c.HTML(code, "index.tmpl", resp)
}

// PostPage .
// @router /archives/:post_id [GET]
func PostPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.PostPage(ctx, c))
}

// CreatePageAPI .
// @router /api/admin/page [POST]
func CreatePageAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.CreatePageAPI(ctx, c))
}

// GetPageAPI .
// @router /api/admin/page/:page_id [GET]
func GetPageAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetPageAPI(ctx, c))
}

// UpdatePageAPI .
// @router /api/admin/page/:page_id [PUT]
func UpdatePageAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.UpdatePageAPI(ctx, c))
}

// GetPageListAPI .
// @router /api/admin/page/list [GET]
func GetPageListAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.GetPageListAPI(ctx, c))
}

// DeletePageAPI .
// @router /api/admin/page/:page_id [DELETE]
func DeletePageAPI(ctx context.Context, c *app.RequestContext) {
	c.JSON(facade.DeletePageAPI(ctx, c))
}

// PagePage .
// @router /pages/:page_slug [GET]
func PagePage(ctx context.Context, c *app.RequestContext) {
	c.HTML(facade.PagePage(ctx, c))
}
