package errorResp

import "google.golang.org/grpc/codes"

var Success = NewApiError(codes.OK, "")