package facade

import (
	"context"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/resp"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func ArchivesPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	rsp, pErr := service.ArchivesPage(ctx)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}
