package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/manmanxing/errors"
	"github.com/manmanxing/go_common/api/errorResp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//针对 error 做一些处理
func HookError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		err := next(context)
		outPutErr(context, err)
		return nil
	}
}

func outPutErr(ctx echo.Context, err error) {
	if err == nil {
		if ctx.Response().Committed {
			return
		}
		e := ctx.JSON(http.StatusOK, errorResp.Success)
		if e != nil {
			e = errors.Wrap(e)
			fmt.Println("outPutErr err", e)
			return
		}
	}

	if ctx.Response().Committed {
		return
	}

	err = errors.Cause(err)
	if _, ok := err.(*errorResp.ApiError); ok {
		err = checkGRPCError(err)
	}
	e := ctx.JSON(http.StatusOK, errorResp.Success)
	if e != nil {
		e = errors.Wrap(e)
		fmt.Println("outPutErr err", e)
		return
	}
}

func checkGRPCError(err error) error {
	if err == nil {
		return errorResp.Success
	}
	//这里判断是不是 grpc 错误
	s, ok := status.FromError(err)
	if !ok {
		return errorResp.NewApiError(codes.Unknown, err.Error())
	}
	if s.Code() == codes.OK {
		//这里屏蔽掉成功时的msg
		return errorResp.Success
	}
	//拆解 grpc 的错误，并组装成统一信息返回
	//todo 可以根据是线上还是测试环境，返回不同的错误提示
	return errorResp.NewApiError(s.Code(), s.Message())
}
