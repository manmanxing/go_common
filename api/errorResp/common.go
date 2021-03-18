package errorResp

import "google.golang.org/grpc/codes"

var Success = NewApiError(codes.OK, "")
var ServerError = NewApiError(codes.Internal, "系统错误，请稍后重试")