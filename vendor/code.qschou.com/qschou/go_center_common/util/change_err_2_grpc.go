package util

import (
	"database/sql"

	"code.qschou.com/golang/alert"

	"code.qschou.com/golang/errors"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ChangeErr2Grpc(err error) error {
	if err == nil {
		return nil
	}
	e := errors.Cause(err)
	if e == nil {
		return nil
	}
	ae, aeok := e.(*alert.Error)
	var pe error
	if aeok {
		if ae != nil {
			pe = ae.BasicError() // TODO 如果是alert.Error，取到BasicError判断，决定返回给客户端什么error
		} else {
			pe = e
		}
	} else {
		pe = e
	}
	if _, ok := status.FromError(pe); ok {
		return pe
	}
	if pe == gorm.ErrRecordNotFound {
		return status.Errorf(codes.NotFound, "Sorry，资源未找到")
	}
	if pe == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "不好意思，资源未找到")
	}
	if pe == redis.ErrNil {
		return status.Errorf(codes.NotFound, "对不起，资源为找到")
	}
	if pe != nil {
		return status.Errorf(codes.Internal, pe.Error())
	}
	return pe
}
