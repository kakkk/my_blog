package facade

import (
	"context"
	"net/http"

	"my_blog/biz/mock"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func ArchivesPage(ctx context.Context, c *app.RequestContext) (resp *page.ArchivesPageResp, code int) {
	return mock.ArchivesPageRespMocker(), http.StatusOK
}
