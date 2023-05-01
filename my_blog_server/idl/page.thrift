namespace go blog.page

const string PageTypeTagList = "PAGE_TAG_LIST"
const string PageTypeCategoryList = "PAGE_CATEGORY_LIST"
const string PageTypeArchives = "PAGE_ARCHIVES"
const string PageTypePost = "PAGE_POST"
const string PageTypeSearch = "PAGE_SEARCH"
const string PageTypeIndex = "PAGE_HOME"
const string PageTypePostList = "PAGE_POST_LIST"
const string PageTypeTagPostList = "PAGE_TAG_POST_LIST"
const string PageTypeCategoryPostList = "PAGE_CATEGORY_POST_LIST"
const string PageTypeError = "PAGE_ERROR"

// ============PageMeta=============

struct PageMeta {
    1:required string Title
    2:required string Description
    3:required string CDNDomain
    4:required string SiteDomain
    5:required string PageType
    6:required string ErrorCode
}

// ==============文章列表=============

struct PostItem {
    1:required string ID
    2:required string Title
    3:required string Abstract
    4:required string Info
}

struct PostListPageRequest {
    1:optional i64 Page (api.path="page")
    2:optional string Name (api.path="name")
    3:optional string PageType
}

struct PostListPageResp {
    1:required string Name
    2:required string Slug
    3:required list<PostItem> PostList
    4:required string PrevPage
    5:required string NextPage

    255:required PageMeta Meta
}

// =============文章归档=============

struct ArchiveByMonth {
    1: required list<PostItem> Posts
    2: required string Month
    3: required string Count
}

struct ArchiveByYear {
    1: required list<ArchiveByMonth> Archives
    2: required string Year
    3: required string Count
}

struct ArchivesPageResp {
    1: required list<ArchiveByYear> PostArchives

    255:required PageMeta Meta
}

// ==============标签云=============

struct TermListItem {
    1:required string Name
    2:required string Count
    3:required string Slug
}

struct TermsPageResp {
    1:required list<TermListItem> List

    255: required PageMeta Meta
}

// ==============基本页面===============

struct BasicPageResp {

    255: required PageMeta Meta
}

// ===============文章页================

struct PostNav {
    1: required string ID
    2: required string Title
}

struct PostPageRequest {
    1: required i64 ID (api.path="post_id")

}

struct PostInfo {
    1: required string Author
    2: required string PublishAt
    3: required string UV
    4: required string WordCount
    5: required list<TermListItem> CategoryList
}

struct PostPageResponse {
    1: required string Title
    2: required PostInfo Info
    3: required string Content
    4: optional list<string> Tags
    5: optional PostNav PrevPage
    6: optional PostNav NextPage

    255: required PageMeta Meta
}


service PageService {
    // =============文章列表==============
    // 首页
    PostListPageResp IndexPage() (api.get="/");
    PostListPageResp IndexByPaginationPage(1: PostListPageRequest request) (api.get="/page/:page");
    // 分类
    PostListPageResp CategoryPostPage(1: PostListPageRequest request) (api.get="/category/:name")
    PostListPageResp CategoryPostByPaginationPage(1: PostListPageRequest request) (api.get="/category/:name/:page")
    // 标签
    PostListPageResp TagPostPage(1: PostListPageRequest request) (api.get="/tag/:name")
    PostListPageResp TagPostByPaginationPage(1: PostListPageRequest request) (api.get="/tag/:name/:page")

    // =============文章归档==============
    ArchivesPageResp ArchivesPage() (api.get="/archives")

    // =============标签云================
    TermsPageResp TagsPage() (api.get="/tags")
    TermsPageResp CategoriesPage() (api.get="/categories")

    // ==============单页面=================
    BasicPageResp SearchPage() (api.get="/search")

    // ==============文章页=================
    PostPageResponse PostPage(1: PostPageRequest request) (api.get="/archives/:post_id")

}