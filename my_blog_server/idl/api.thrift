namespace go blog.api


struct BaseResp {
    1: required i32 status_code
    2: required string status_message
}

// 登录
struct LoginRequest{
    1: required string user_name
    2: required string password
}

struct LoginResponse{

    255:required BaseResp base_resp (go.tag="json:\"-\"")
}


service APIService {
    LoginResponse LoginAPI(1:LoginRequest request) (api.post="/api/login")
}