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

type IAPIResponse interface {
	GetBaseResp() (v *api.BaseResp)
}

func NewBaseResponse(code common.RespCode, msg string) *api.BaseResp {
	return &api.BaseResp{
		Code: code,
		Msg:  msg,
	}
}

func NewAPIResponse(resp IAPIResponse) *APIResponse {
	if resp.GetBaseResp().GetCode() != common.RespCode_Success {
		return &APIResponse{
			Code: resp.GetBaseResp().GetCode(),
			Msg:  resp.GetBaseResp().GetMsg(),
			Data: nil,
		}
	}
	return &APIResponse{
		Code: resp.GetBaseResp().GetCode(),
		Msg:  resp.GetBaseResp().GetMsg(),
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

func NewFailResp() *APIResponse {
	return &APIResponse{
		Code: common.RespCode_Fail,
		Msg:  "fail",
		Data: nil,
	}
}

func NewFailBaseResp() *api.BaseResp {
	return NewBaseResponse(common.RespCode_Fail, "fail")
}

func NewSuccessBaseResp() *api.BaseResp {
	return NewBaseResponse(common.RespCode_Success, "ok")
}
