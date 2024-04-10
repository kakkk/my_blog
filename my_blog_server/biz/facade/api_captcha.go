package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/api"
)

func GetCaptchaAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var rsp api.GetCaptchaAPIResponse
	return http.StatusOK, resp.NewAPIResponse(&rsp)
}
