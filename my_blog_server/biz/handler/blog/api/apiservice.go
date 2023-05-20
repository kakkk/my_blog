// Code generated by hertz generator.

package api

import (
	"context"

	"my_blog/biz/facade"
	_ "my_blog/biz/model/blog/api"

	"github.com/cloudwego/hertz/pkg/app"
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
