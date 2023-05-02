package facade

import (
	"context"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/errorx"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func PostPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostPageRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("parameter error:[%v]", err)
		return resp.PackPageResponse(nil, errorx.NewNotFoundErrPageError(), consts.IndexTmpl)
	}
	rsp, pErr := service.PostPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}
