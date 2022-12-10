namespace go blog.common


// =========错误码=========

enum RespCode {
    Success         = 0
    Fail            = 1
    LoginFail       = 200001
    ParameterError  = 400000
    Unauthorized    = 400100
    InternalError   = 500000
}

// =========Extra信息===========

enum ExtraInfo {
    CategoryOrder = 1 // 分类排序
}

//==========文章相关============

enum ArticleType {
    Post = 1    // 博客文章
    Page = 2    // 独立页面
}

enum ArticleStatus {
    DRAFT   = 1     // 草稿
    PUBLIST = 2     // 发布
    OFFLINE = 3     // 下线
    DELETE  = 4     // 删除
}