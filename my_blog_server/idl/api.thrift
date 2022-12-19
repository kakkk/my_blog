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
}