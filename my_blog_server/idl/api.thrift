namespace go blog.api

struct TestReq {
    1: string Name (api.json="name");
}

struct TestResp {
    1: string RespBody;
}


service APIService {
    TestResp HelloMethod(1: TestReq request) (api.get="/test");
}