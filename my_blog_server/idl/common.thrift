namespace go blog.common


// ======DeleteFlag=======
enum DeleteFlag {
    Exist   = 0     // 存在
    Delete  = 1     // 删除
}

// =========错误码=========
enum RespCode {
    Success         = 0         // 成功
    Fail            = 1         // 失败
    LoginFail       = 200001    // 登录失败
    HasExist        = 200002    // 已存在
    NotFound        = 200003    // 不存在
    ParameterError  = 400000    // 参数错误
    Unauthorized    = 400100    // 未登录
    InternalError   = 500000    // 服务错误
}

// =========Extra信息===========
enum ExtraInfo {
    CategoryOrder   = 1     // 分类排序
    DefaultCategory = 2     // 默认分类
}

//==========文章相关============

enum ArticleType {
    Post = 1    // 博客文章
    Page = 2    // 独立页面
}

enum ArticleStatus {
    DRAFT   = 1     // 草稿
    PUBLISH = 2     // 发布
    OFFLINE = 3     // 下线
    DELETE  = 4     // 删除
}