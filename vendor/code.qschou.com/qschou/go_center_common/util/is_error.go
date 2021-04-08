package util

import (
	"code.qschou.com/golang/alert"
	"code.qschou.com/golang/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.Cause(err)
	if e == nil {
		return false
	}
	_, aeok := e.(*alert.Error)
	if aeok { // TODO 如果是alert.Error必是error，必打error日志
		return true
	}
	if s, ok := status.FromError(e); ok {
		if s.Code() > 100 {
			return false
		}
		switch s.Code() {
		case codes.InvalidArgument, codes.Unauthenticated, codes.FailedPrecondition, codes.Unavailable:
			return false
		default:
			return true
		}
	}
	return true
}
