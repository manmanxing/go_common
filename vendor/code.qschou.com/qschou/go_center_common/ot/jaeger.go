package ot

import (
	"net/http"
	"strings"
	"time"

	"code.qschou.com/qschou/go_center_common/dlog"

	"code.qschou.com/qschou/go_center_common/gls"
	"code.qschou.com/qschou/go_center_common/util"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"

	"golang.org/x/net/context"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
)

//jaeger中间件，用于对每个请求生成一个span，并使用opentracing生成的spanID作为日志的spanID
func JaegerTrace(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		//find spanContext in http header
		parentSpanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request().Header))
		optName := c.Request().RequestURI
		if c.Request().Method == http.MethodGet {
			pos := strings.Index(optName, "?")
			if pos != -1 {
				optName = optName[:pos]
			}
		}
		span := opentracing.StartSpan(optName,
			ext.RPCServerOption(parentSpanCtx),
			ext.SpanKindRPCServer,
			opentracing.Tag{Key: string(ext.PeerAddress), Value: c.RealIP()},
			opentracing.Tag{Key: string(ext.HTTPUrl), Value: c.Request().RequestURI},
		)
		defer func() {
			if err != nil {
				span.SetTag("error", err.Error())
			}
			span.Finish()
		}()
		ctx := context.Background()
		ctx = opentracing.ContextWithSpan(ctx, span) //这里context关联span
		//SetGlobalContext(c, ctx)
		//获取tracerID、spanID和pSpanID
		traceID := getTracerID(c, span.Context())
		pSpanID := getParentSpanID(c, parentSpanCtx, span.Context())
		spanID := getSpanID(c, span.Context())

		gls.SetGlsWithCtx(traceID, pSpanID, spanID, ctx, func() {
			err = next(c)
		})
		return err
	}

}

func getTracerID(c echo.Context, spanCtx opentracing.SpanContext) (tracerID string) {
	tracerIDTmp, ok := spanCtx.(jaeger.SpanContext)
	if !ok {
		return
	}
	return tracerIDTmp.TraceID().String()
	//tracerID = c.Request().Header.Get("qsc-header-tid")
	return
}

func getParentSpanID(c echo.Context, parentSpanCtx, spanCtx opentracing.SpanContext) (pSpanID string) {
	//先从opentracing中取parentSpan的spanid，作为本次请求的 pSpanID
	_, pSpanID, _ = GetSpanIDFromSpanCtx(parentSpanCtx)
	if pSpanID != "" && pSpanID != "0" {
		return
	}
	//没有的话，取Header中的spanID
	pSpanID = c.Request().Header.Get("qsc-header-spanid")
	if pSpanID != "" {
		return
	}
	//还没有的话，pSpanID == spanID
	_, pSpanID, _ = GetSpanIDFromSpanCtx(spanCtx)
	return
}

func getSpanID(c echo.Context, spanCtx opentracing.SpanContext) (spanID string) {
	_, spanID, _ = GetSpanIDFromSpanCtx(spanCtx)
	if spanID != "" {
		return
	}
	spanID = util.GenerateSpanID(c.Request().RemoteAddr)
	return
}

const (
	globalContextKey = `global_context`
)

//将context设置到echo的context中
func SetGlobalContext(c echo.Context, ctx context.Context) {
	c.Set(globalContextKey, ctx)
}

func GetGlobalContext(c echo.Context) context.Context {
	cc := c.Get(globalContextKey)
	if cc == nil {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		return ctx
	}
	return cc.(context.Context)
}

// 从spanContext中获取父spanid和spanid，目前仅支持jaeger的
func GetSpanIDFromSpanCtx(spanCtx opentracing.SpanContext) (pSpanID, spanID string, err error) {
	switch t := spanCtx.(type) {
	case jaeger.SpanContext:
		spanID = spanCtx.(jaeger.SpanContext).SpanID().String()
		pSpanID = spanCtx.(jaeger.SpanContext).ParentID().String()
	default:
		err = errors.Wrapf(errors.New("sorry, unsupported spanCtx type."), "type", t)
		dlog.Warn("GetSpanIDFromSpanCtx", err)
	}
	return
}
