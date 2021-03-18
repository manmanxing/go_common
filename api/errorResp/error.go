package errorResp

import (
	"encoding/json"
	"fmt"
	"github.com/manmanxing/errors"
	"google.golang.org/grpc/codes"
)

//统一 error 信息输出
//当没有错误时，返回的code 为0，msg 为""
//这里基本指的是 grpc 的错误
type ApiError struct {
	Code codes.Code `json:"code"`
	Msg  string     `json:"msg"`
}

func (ae *ApiError) Error() string {
	bt, err := json.Marshal(ae)
	if err != nil {
		err = errors.Wrap(err)
		fmt.Println("apiError json marshal err", err)
		return ""
	}
	return string(bt)
}

func NewApiError(code codes.Code, msg string) *ApiError {
	return &ApiError{
		Code: code,
		Msg:  msg,
	}
}

//统一信息返回
type StdInfoResp struct {
	*ApiError
	Data interface{} `json:"data"`
}

func (sir *StdInfoResp) Error() string {
	bt, err := json.Marshal(sir)
	if err != nil {
		err = errors.Wrap(err)
		fmt.Println("stdInfoResp json marshal err", err)
		return ""
	}
	return string(bt)
}

//这里可以将返回的信息进行分类
func NewStdInfoResp(data interface{}, err error) *StdInfoResp {
	resp := new(StdInfoResp)
	resp.Data = data
	if err == nil {
		resp.ApiError = Success
		return resp
	}
	//如果是开发定义的 apiError 类型，就直接返回
	if val, ok := err.(*ApiError); ok {
		resp.ApiError = val
		return resp
	} else {
		//如果不是 apiError 类型，组装成 apiError 类型
		//todo 需要针对线上环境屏蔽掉错误详情？
		resp.ApiError = NewApiError(codes.Unknown, err.Error())
	}
	return resp
}