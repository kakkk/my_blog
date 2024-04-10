package facade

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/consts"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/page"
)

func PostPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostPageRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("parameter error:[%v]", err)
		return resp.PackPageResponse(nil, errorx.NewNotFoundErrPageError(), consts.IndexTmpl)
	}
	rsp, pErr := application.GetBlogApplication().PostPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}
