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

// ========= 独立页面 =========

// 创建独立页面
struct CreatePageAPIRequest {
    1: required string Title
    2: required string Content
    3: required string Slug

}

// 创建独立页面
struct CreatePageAPIResponse {
    1: required i64 ID

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

// 获取独立页面
struct GetPageAPIRequest {
    1: required i64 ID (api.path="page_id")
}

// 获取独立页面
struct GetPageAPIResponse {
    1: required i64 ID
    2: required string Title
    3: required string Content
    4: required string Slug

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}

struct PageListItem {
    1: required i64 ID
    2: required string Title
    3: required string Slug

}

// 更新页面
struct UpdatePageAPIRequest {
    1: required i64 ID (api.path="page_id")
    2: required string Title
    3: required string Content
    4: required string Slug

}

// 删除页面
struct DeletePageAPIRequest {
    1: required i64 ID (api.path="page_id")

}

// 获取独立页面列表
struct GetPageListAPIResponse {
    1: optional list<PageListItem> PageList

    255: required BaseResp BaseResp (go.tag="json:\"-\"")
}