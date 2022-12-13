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
