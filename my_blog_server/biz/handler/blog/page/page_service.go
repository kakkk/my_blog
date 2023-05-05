// Code generated by hertz generator.

package page

import (
	"context"

	"my_blog/biz/facade"
	_ "my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

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
	resp, code := facade.CategoryPostPage(ctx, c)
	c.HTML(code, "index.tmpl", resp)
}

// CategoryPostByPaginationPage .
// @router /category/:slug/:page [GET]
func CategoryPostByPaginationPage(ctx context.Context, c *app.RequestContext) {
	resp, code := facade.CategoryPostByPaginationPage(ctx, c)
	c.HTML(code, "index.tmpl", resp)
}

// TagPostPage .
// @router /tag/:name [GET]
func TagPostPage(ctx context.Context, c *app.RequestContext) {
	resp, code := facade.TagPostPage(ctx, c)
	c.HTML(code, "index.tmpl", resp)
}

// TagPostByPaginationPage .
// @router /tag/:name/:page [GET]
func TagPostByPaginationPage(ctx context.Context, c *app.RequestContext) {
	resp, code := facade.TagPostByPaginationPage(ctx, c)
	c.HTML(code, "index.tmpl", resp)
}

// ArchivesPage .
// @router /archives [GET]
func ArchivesPage(ctx context.Context, c *app.RequestContext) {
	resp, code := facade.ArchivesPage(ctx, c)
	c.HTML(code, "index.tmpl", resp)
}

// TagsPage .
// @router /tags [GET]
func TagsPage(ctx context.Context, c *app.RequestContext) {
	resp, code := facade.TagsPage(ctx, c)
	c.HTML(code, "index.tmpl", resp)
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
