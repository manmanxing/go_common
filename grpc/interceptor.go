package grpc

import (
	"code.qschou.com/qschou/go_center_common/ot"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/manmanxing/go_common/gls"
	"github.com/manmanxing/go_common/util"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

/*
	设置 grpc 拦截器
*/

//这个中间件统一注入opentracing的span
func unaryClientInterceptForInjectSpan(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	sp := opentracing.SpanFromContext(ctx)
	if sp == nil {
		ctx = opentracing.ContextWithSpan(ctx, opentracing.SpanFromContext(gls.TraceCtx()))
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

//针对一元模式的客户端请求打印日志
func unaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	//这里若没有设置超时，就强制设置
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
	}

	p := peer.Peer{}
	if opts == nil {
		opts = []grpc.CallOption{grpc.Peer(&p)}
	} else {
		opts = append(opts, grpc.Peer(&p))
	}

	start := time.Now()
	defer func() {
		in, _ := json.Marshal(req)
		out, _ := json.Marshal(reply)
		inStr, outStr := string(in), string(out)
		dur := int64(time.Since(start) / time.Millisecond)
		var serviceIp string
		if p.Addr != nil {
			serviceIp = p.Addr.String()
		}
		if dur >= 300 {
			log.Warn("grpc", method, "in", inStr, "out", outStr, "err", err, "dur/ms", dur, "ip", serviceIp, "type", "access_result")
		} else {
			log.Info("grpc", method, "in", inStr, "out", outStr, "err", err, "dur/ms", dur, "ip", serviceIp, "type", "access_result")
		}
	}()
	return invoker(ctx, method, req, reply, cc, opts...)
}

//针对一元模式的服务端请求打印日志
func unaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	traceID, pSpanID, spanID := ot.TraceInfo(ctx)
	gls.SetGlsWithCtx(traceID, pSpanID, spanID, ctx, func() {
		resp, err = _UnaryServerInterceptor(ctx, req, info, handler)
	})
	return
}

func _UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	remote, _ := peer.FromContext(ctx)
	remoteAddr := remote.Addr.String()

	bt, _ := json.Marshal(req)
	log.Info("ip", remoteAddr, "grpc_access_start", info.FullMethod, "in", string(bt))

	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			var recv error
			switch e := r.(type) {
			case error:
				recv = e
			default:
				recv = fmt.Errorf("%v", e)
			}

			stackInfo := util.StackInfo(util.Callers(3))
			log.Error("panic", recv, "stack", stackInfo)
			err = status.Error(codes.Internal, fmt.Sprintf("panic=%s", recv))
		}

		respBt, _ := json.Marshal(resp)
		dur := int64(time.Since(start) / time.Millisecond)
		if dur >= 300 {
			log.Warn("ip", remoteAddr, "grpc_access_end", info.FullMethod, "in", string(bt), "out", string(respBt), "err", err, "dur/ms", dur)
		} else {
			log.Info("ip", remoteAddr, "grpc_access_end", info.FullMethod, "in", string(bt), "out", string(respBt), "err", err, "dur/ms", dur)
		}
	}()

	resp, err = handler(ctx, req)
	if err != nil {
		_, ok := status.FromError(err)
		if !ok {
			err = status.Error(codes.Internal, err.Error())
		}
	}

	return
}

//todo 针对流rpc模式打印日志