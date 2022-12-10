package resp

import (
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/common"
)

type APIResponse struct {
	Code common.RespCode `json:"code"`
	Msg  string          `json:"msg"`
	Data any             `json:"data"`
}

func NewBaseResponse(code common.RespCode, msg string) *api.BaseResp {
	return &api.BaseResp{
		Code: code,
		Msg:  msg,
	}
}

func NewAPIResponse(base *api.BaseResp, resp any) *APIResponse {
	if base.GetCode() != common.RespCode_Success {
		return &APIResponse{
			Code: base.Code,
			Msg:  base.Msg,
			Data: nil,
		}
	}
	return &APIResponse{
		Code: base.Code,
		Msg:  base.Msg,
		Data: resp,
	}
}

func NewInternalErrorResp() *APIResponse {
	return &APIResponse{
		Code: common.RespCode_InternalError,
		Msg:  "Internal Error",
		Data: nil,
	}
}

func NewParameterErrorResp() *APIResponse {
	return &APIResponse{
		Code: common.RespCode_ParameterError,
		Msg:  "Parameter Error",
		Data: nil,
	}
}

func NewFailBaseResp() *api.BaseResp {
	return NewBaseResponse(common.RespCode_Fail, "fail")
}
