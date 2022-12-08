// Code generated by hertz generator.

package api

import (
	"context"

	serviceResp "my_blog/biz/common/resp"
	"my_blog/biz/model/blog/api"

	"github.com/cloudwego/hertz/pkg/app"
)

// LoginAPI .
// @router /api/login [POST]
func LoginAPI(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	resp := &api.LoginResponse{
		TestMsg:  "test",
		BaseResp: serviceResp.GetBaseResp(0, ""),
	}

	c.JSON(200, serviceResp.GetAPIResponse(resp.GetBaseResp(), resp))
}
