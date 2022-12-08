package resp

import "my_blog/biz/model/blog/api"

type APIResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func GetBaseResp(code int32, msg string) *api.BaseResp {
	return &api.BaseResp{
		StatusCode:    code,
		StatusMessage: msg,
	}
}

func GetAPIResponse(base *api.BaseResp, resp any) *APIResponse {
	return &APIResponse{
		Code: base.StatusCode,
		Msg:  base.StatusMessage,
		Data: resp,
	}
}
