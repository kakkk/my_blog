namespace go blog.api


struct BaseResp {
    1: required i32 StatusCode
    2: required string StatusMessage
}

// 登录
struct LoginRequest{
    1: required string Username
    2: required string Password
}

struct LoginResponse{
    1: required i64 UserID
    2: required string Username
    3: required string Nickname

    255:required BaseResp BaseResp (go.tag="json:\"-\"")
}


service APIService {
    LoginResponse LoginAPI(1:LoginRequest request) (api.post="/api/admin/login")
}