namespace go blog.api

include "common.thrift"

// BaseResp
struct BaseResp {
    1: required common.RespCode Code
    2: required string Msg
}

// 通用Response
struct CommonResponse {

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 分页
struct Pagination {
    1: required i32 Page
    2: required i32 Limit
    3: required bool HasMore
    4: optional i64 Total
}

// ===========用户相关==========

// 登录
struct LoginRequest {
    1: required string Username
    2: required string Password
}

// 登录
struct LoginResponse {
    1: required i64 UserID
    2: required string Username
    3: required string Nickname
    4: required string Avatar

    255:required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 获取用户列表
struct GetUserInfoAPIResponse {
    1: required i64 UserID
    2: required string Username
    3: required string Nickname
    4: required string Avatar

    255:required BaseResp BaseResp (go.tag="json:\"-\"")
}

// ==========标签相关===========

// 创建标签
struct CreateTagAPIRequest {
    1: required string Name
}

// 创建标签
struct CreateTagAPIResponse {
    1: required i64 ID
    2: required string Name

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 更新标签
struct UpdateTagAPIRequest{
    1: required i64 ID (api.path="tag_id")
    2: required string Name
}

// 删除标签
struct DeleteTagAPIRequest {
    1: required i64 ID (api.path="tag_id")
}

// 获取标签列表
struct GetTagListAPIRequest {
    1: optional string Keyword (api.query="keyword")
    2: optional i32 Page (api.query="page")
    3: optional i32 Limit (api.query="limit")
}

// 获取标签列表
struct TagListItem {
    1: required i64 ID
    2: required string Name
    3: required i64 Count
}

// 获取标签列表
struct GetTagListAPIResponse {
    1: required Pagination Pagination
    2: required list<TagListItem> TagList

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// ==========分类相关=============

// 创建分类
struct CreateCategoryAPIRequest{
    1: required string Name
    2: required string Slug
}

// 创建分类
struct UpdateCategoryAPIRequest{
    1: required i64 ID (api.path="category_id")
    2: required string Name
    3: required string Slug
}

// 更新分类排序
struct UpdateCategoryOrderAPIRequest {
    1: required list<i64> Order
}

// 删除分类
struct DeleteCategoryAPIRequest {
    1: required i64 ID (api.path="category_id")
}

// 分类
struct CategoryListItem {
    1: required i64 ID
    2: required string Name
    3: required string Slug
    4: required i64 Count

}

// 获取分类列表
struct GetCategoryListAPIResponse {
    1: required list<CategoryListItem> CategoryList

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}


// ==========文章相关=============

// 创建文章
struct CreatePostAPIRequest {
    1: required string Title
    2: required string Content
    3: required common.ArticleStatus Status
    4: required list<i64> CategoryList
    5: required list<string> Tags

}

// 创建文章
struct CreatePostAPIResponse {
    1: required i64 ID

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 获取文章
struct GetPostAPIRequest {
    1: required i64 ID (api.path="post_id")
}

// 获取文章 - 分类
struct CategoriesItem {
    1: required i64 ID
    2: required string Name
}

// 获取文章
struct GetPostAPIResponse {
    1: required i64 ID
    2: required string Title
    3: required string Content
    4: required common.ArticleStatus Status
    5: required list<CategoriesItem> CategoryList
    6: required list<string> Tags

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 更新文章
struct UpdatePostAPIRequest {
    1: required i64 ID (api.path="post_id")
    2: required string Title
    3: required string Content
    4: required list<i64> CategoryList
    5: required list<string> Tags

}

// 更新文章状态
struct UpdatePostStatusAPIRequest {
    1: required i64 ID (api.path="post_id")
    2: required common.ArticleStatus Status
}

// 获取文章列表
struct GetPostListAPIRequest {
    1: optional string Keyword (api.query="keyword")
    2: optional list<string> Categories (api.query="categories")
    3: optional list<string> Tags (api.query="tags")
    4: optional i32 Page (api.query="page")
    5: optional i32 Limit (api.query="limit")
}

// 文章列表
struct PostListItem {
    1: required i64 ID
    2: required string Title
    3: required list<string> CategoryList
    4: required string Editor
    5: required common.ArticleStatus Status
    6: required i64 UV
    7: required i64 UpdateAt
    8: required i64 PublishAt

}

// 获取文章列表
struct GetPostListAPIResponse {
    1: required Pagination Pagination
    2: optional list<PostListItem> PostList

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 删除文章
struct DeletePostAPIRequest {
    1: required i64 ID (api.path="post_id")
}

// 获取评论列表
struct GetCommentListAdminAPIRequest {
    1: optional string Email
    2: optional string Nickname
    3: optional i64 ArticleID
    4: optional list<common.CommentStatus> Status
    5: optional i32 Page (api.query="page")
    6: optional i32 Limit (api.query="limit")
}

struct ArticleMeta {
    1: required i64 ID
    2: required string Title
}

struct GetCommentListAdminItem {
    1: required i64 ID
    2: required string Nickname
    3: required string Avatar
    4: required string Website
    5: required ArticleMeta Article
    6: required string Content
    7: optional i64 ReplyToID
    8: optional string ReplyToContent
    9: required common.CommentStatus Status
}

// 获取评论列表
struct GetCommentListAdminAPIResponse {
    1: required Pagination Pagination
    2: optional list<PostListItem> PostList

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

struct ReplyCommentAdminAPIRequest {
    1: required i64 CommentID (api.path="comment_id")
    2: required string Content
}

struct UpdateCommentStatusAdminAPIRequest {
    1: required i64 CommentID (api.path="comment_id")
    2: required common.CommentStatus Status
}

struct DeleteCommentAdminAPIRequest {
    1: required i64 CommentID (api.path="comment_id")
}

// ==========用户侧接口=============

// 搜索
struct SearchAPIRequest {
    1: required string Query (api.query="q")
}

struct SearchResultItem {
    1: required string Link
    2: required string Title
    3: required string Abstract
}

struct SearchAPIResponse {
    1: optional list<SearchResultItem> Results

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 评论
struct Comment {
    1: required i64 ID
    2: required string Nickname
    3: required string Avatar
    4: required string Website
    5: required string Content
    6: required string CommentAt
    7: required string ReplyUser
}

struct CommentListItem {
    1: required Comment Comment
    2: optional list<Comment> Replies
}

struct GetCommentListAPIResponse {
    1: required list<CommentListItem> Comments
    2: required bool HasMore

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

struct GetCommentListAPIRequest {
    1: required i64 ArticleID (api.query="article_id")
}

struct CommentArticleAPIRequest {
    1: required i64 ArticleID           // 文章ID
    2: required string Nickname         // 昵称
    3: required string Email            // 邮箱
    4: optional string Website          // 网址
    5: required string Content          // 评论内容
    6: required string VerifyID         // 验证码ID
    7: required string VerifyCode       // 验证码
}

struct CommentArticleAPIResponse {
    1: required i64 ID
    2: required common.CommentStatus CommentStatus

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

struct ReplyCommentAPIRequest {
    1: required i64 ArticleID           // 文章ID
    2: optional i64 ReplyID             // 回复的评论ID
    3: required string Nickname         // 昵称
    4: required string Email            // 邮箱
    5: optional string Website          // 网址
    6: required string Content          // 评论内容
    7: required string VerifyID         // 验证码ID
    8: required string VerifyCode       // 验证码
}

struct ReplyCommentAPIResponse {
    1: required i64 ID
    2: required common.CommentStatus CommentStatus

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 验证码
struct GetCaptchaAPIResponse {
    1: required string Captcha

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}



service APIService {
    // ==========用户相关=============
    // 登录
    LoginResponse LoginAPI(1:LoginRequest request) (api.post="/api/admin/login")
    // 获取用户信息
    GetUserInfoAPIResponse GetUserInfoAPI() (api.get="/api/admin/user/info")
    // ==========标签相关=============
    // 创建标签
    CreateTagAPIResponse CreateTagAPI(1:CreateTagAPIRequest request) (api.post="/api/admin/tag")
    // 更新标签
    CommonResponse UpdateTagAPI(1:UpdateTagAPIRequest request) (api.put="/api/admin/tag/:tag_id")
    // 删除标签
    CommonResponse DeleteTagAPI(1:DeleteTagAPIRequest request) (api.delete="/api/admin/tag/:tag_id")
    // 获取标签列表
    GetTagListAPIResponse GetTagListAPI(GetTagListAPIRequest request) (api.get="/api/admin/tag/list")
    // ==========分类相关=============
    // 创建分类
    CommonResponse CreateCategoryAPI(1:CreateCategoryAPIRequest request) (api.post="/api/admin/category")
    // 更新分类
    CommonResponse UpdateCategoryAPI(1:UpdateCategoryAPIRequest request) (api.put="/api/admin/category/:category_id")
    // 删除分类
    CommonResponse DeleteCategoryAPI(1:DeleteCategoryAPIRequest request) (api.delete="/api/admin/category/:category_id")
    // 更新分类排序
    CommonResponse UpdateCategoryOrderAPI(1:UpdateCategoryOrderAPIRequest request) (api.put="/api/admin/category/order")
    // 获取分类列表
    GetCategoryListAPIResponse GetCategoryListAPI() (api.get="/api/admin/category/list")
    // ==========文章相关=============
    // 创建文章
    CreatePostAPIResponse CreatePostAPI(1:CreatePostAPIRequest request) (api.post="/api/admin/post")
    // 获取文章
    GetPostAPIResponse GetPostAPI(1:GetPostAPIRequest request) (api.get="/api/admin/post/:post_id")
    // 更新文章
    CommonResponse UpdatePostAPI(1:UpdatePostAPIRequest request) (api.put="/api/admin/post/:post_id")
    // 更新文章状态
    CommonResponse UpdatePostStatusAPI(1:UpdatePostStatusAPIRequest request) (api.put="/api/admin/post/:post_id/status")
    // 获取文章列表
    GetPostListAPIResponse GetPostListAPI(1:GetPostListAPIRequest request) (api.post="/api/admin/post/list")
    // 删除文章
    CommonResponse DeletePostAPI(1:DeletePostAPIRequest request) (api.delete="/api/admin/post/:post_id")
    // ==========评论相关=============
    // 获取评论列表
    GetCommentListAdminAPIResponse GetCommentListAdminAPI(1:GetCommentListAdminAPIRequest request) (api.get="/api/admin/comment/list")
    // 管理员回复评论
    CommonResponse ReplyCommentAdminAPI(1:ReplyCommentAdminAPIRequest request) (api.post="/api/admin/comment/:comment_id/reply")
    // 修改评论状态
    CommonResponse UpdateCommentStatusAdminAPI(1:UpdateCommentStatusAdminAPIRequest request) (api.put="/api/admin/comment/:comment_id/status")
    // 删除评论
    CommonResponse DeleteCommentAdminAPI(1:DeleteCommentAdminAPIRequest request) (api.put="/api/admin/comment/:comment_id")

    // ==========用户侧接口============
    // 搜索
    SearchAPIResponse SearchAPI(1:SearchAPIRequest request) (api.get="/api/search")
    // 获取评论列表
    GetCommentListAPIResponse GetCommentListAPI(1:GetCommentListAPIRequest request) (api.get="/api/comment/list")
    // 评论文章
    CommentArticleAPIResponse CommentArticleAPI(1:CommentArticleAPIRequest requset) (api.post="/api/comment/article")
    // 回复评论
    ReplyCommentAPIResponse ReplyCommentAPI(1:ReplyCommentAPIRequest requset) (api.post="/api/comment/reply")
    // 获取验证码
    GetCaptchaAPIResponse GetCaptchaAPI() (api.get="/api/captcha")

}