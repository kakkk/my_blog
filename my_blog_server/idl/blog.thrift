namespace go blog

include "common.thrift"
include "api.thrift"
include "page.thrift"

service BlogService {
    // =============================================================== API ===============================================================

    // ==========用户相关=============

    // 登录
    api.LoginResponse LoginAPI(1:api.LoginRequest request) (api.post="/api/admin/login")
    // 获取用户信息
    api.GetUserInfoAPIResponse GetUserInfoAPI() (api.get="/api/admin/user/info")

    // ==========标签相关=============

    // 创建标签
    api.CreateTagAPIResponse CreateTagAPI(1:api.CreateTagAPIRequest request) (api.post="/api/admin/tag")
    // 更新标签
    api.CommonResponse UpdateTagAPI(1:api.UpdateTagAPIRequest request) (api.put="/api/admin/tag/:tag_id")
    // 删除标签
    api.CommonResponse DeleteTagAPI(1:api.DeleteTagAPIRequest request) (api.delete="/api/admin/tag/:tag_id")
    // 获取标签列表
    api.GetTagListAPIResponse GetTagListAPI(1:api.GetTagListAPIRequest request) (api.get="/api/admin/tag/list")

    // ==========分类相关=============

    // 创建分类
    api.CommonResponse CreateCategoryAPI(1:api.CreateCategoryAPIRequest request) (api.post="/api/admin/category")
    // 更新分类
    api.CommonResponse UpdateCategoryAPI(1:api.UpdateCategoryAPIRequest request) (api.put="/api/admin/category/:category_id")
    // 删除分类
    api.CommonResponse DeleteCategoryAPI(1:api.DeleteCategoryAPIRequest request) (api.delete="/api/admin/category/:category_id")
    // 更新分类排序
    api.CommonResponse UpdateCategoryOrderAPI(1:api.UpdateCategoryOrderAPIRequest request) (api.put="/api/admin/category/order")
    // 获取分类列表
    api.GetCategoryListAPIResponse GetCategoryListAPI() (api.get="/api/admin/category/list")

    // ==========文章相关=============

    // 创建文章
    api.CreatePostAPIResponse CreatePostAPI(1:api.CreatePostAPIRequest request) (api.post="/api/admin/post")
    // 获取文章
    api.GetPostAPIResponse GetPostAPI(1:api.GetPostAPIRequest request) (api.get="/api/admin/post/:post_id")
    // 更新文章
    api.CommonResponse UpdatePostAPI(1:api.UpdatePostAPIRequest request) (api.put="/api/admin/post/:post_id")
    // 更新文章状态
    api.CommonResponse UpdatePostStatusAPI(1:api.UpdatePostStatusAPIRequest request) (api.put="/api/admin/post/:post_id/status")
    // 获取文章列表
    api.GetPostListAPIResponse GetPostListAPI(1:api.GetPostListAPIRequest request) (api.post="/api/admin/post/list")
    // 删除文章
    api.CommonResponse DeletePostAPI(1:api.DeletePostAPIRequest request) (api.delete="/api/admin/post/:post_id")

    // ==========页面相关=============

    // 创建页面
    api.CreatePageAPIResponse CreatePageAPI(1:api.CreatePageAPIRequest request) (api.post="/api/admin/page")
    // 获取页面
    api.GetPageAPIResponse GetPageAPI(1:api.GetPageAPIRequest request) (api.get="/api/admin/page/:page_id")
    // 更新页面
    api.CommonResponse UpdatePageAPI(1:api.UpdatePageAPIRequest request) (api.put="/api/admin/page/:page_id")
    // 获取页面列表
    api.GetPageListAPIResponse GetPageListAPI() (api.get="/api/admin/page/list")
    // 删除页面
    api.CommonResponse DeletePageAPI(1:api.DeletePageAPIRequest request) (api.delete="/api/admin/page/:page_id")

    // ==========评论相关=============

    // 获取评论列表
    api.GetCommentListAdminAPIResponse GetCommentListAdminAPI(1:api.GetCommentListAdminAPIRequest request) (api.get="/api/admin/comment/list")
    // 管理员回复评论
    api.CommonResponse ReplyCommentAdminAPI(1:api.ReplyCommentAdminAPIRequest request) (api.post="/api/admin/comment/:comment_id/reply")
    // 修改评论状态
    api.CommonResponse UpdateCommentStatusAdminAPI(1:api.UpdateCommentStatusAdminAPIRequest request) (api.put="/api/admin/comment/:comment_id/status")
    // 删除评论
    api.CommonResponse DeleteCommentAdminAPI(1:api.DeleteCommentAdminAPIRequest request) (api.put="/api/admin/comment/:comment_id")

    // ==========用户侧接口============

    // 搜索
    api.SearchAPIResponse SearchAPI(1:api.SearchAPIRequest request) (api.get="/api/search")
    // 获取评论列表
    api.GetCommentListAPIResponse GetCommentListAPI(1:api.GetCommentListAPIRequest request) (api.get="/api/comment/list")
    // 评论文章
    api.CommentArticleAPIResponse CommentArticleAPI(1:api.CommentArticleAPIRequest requset) (api.post="/api/comment/article")
    // 回复评论
    api.ReplyCommentAPIResponse ReplyCommentAPI(1:api.ReplyCommentAPIRequest requset) (api.post="/api/comment/reply")
    // 获取验证码
    api.GetCaptchaAPIResponse GetCaptchaAPI() (api.get="/api/captcha")

     // =============================================================== Page ===============================================================

    // =============文章列表==============

    // 首页
    page.PostListPageResp IndexPage() (api.get="/");
    page.PostListPageResp IndexByPaginationPage(1: page.PostListPageRequest request) (api.get="/page/:page");
    // 分类
    page.PostListPageResp CategoryPostPage(1: page.PostListPageRequest request) (api.get="/category/:name")
    page.PostListPageResp CategoryPostByPaginationPage(1: page.PostListPageRequest request) (api.get="/category/:name/:page")
    // 标签
    page.PostListPageResp TagPostPage(1: page.PostListPageRequest request) (api.get="/tag/:name")
    page.PostListPageResp TagPostByPaginationPage(1: page.PostListPageRequest request) (api.get="/tag/:name/:page")

    // =============文章归档==============

    page.ArchivesPageResp ArchivesPage() (api.get="/archives")

    // =============标签云================

    page.TermsPageResp TagsPage() (api.get="/tags")
    page.TermsPageResp CategoriesPage() (api.get="/categories")

    // ==============单页面=================

    page.BasicPageResp SearchPage() (api.get="/search")

    // ==============文章页=================

    page.PostPageResponse PostPage(1: page.PostPageRequest request) (api.get="/archives/:post_id")
    page.PagePageResponse PagePage(1: page.PagePageRequest request) (api.get="/pages/:page_slug")
}