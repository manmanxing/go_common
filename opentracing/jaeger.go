package opentracing

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/manmanxing/go_common/gls"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

//从span中获取trace信息
func getTraceInfoFromSpan(sp opentracing.Span) (traceId, pSpanId, spanId string) {
	if sp == nil {
		return "", "", ""
	}
	spanContext, ok := sp.Context().(jaeger.SpanContext)
	if !ok {
		fmt.Println("")
		return "", "", ""
	}
	return spanContext.TraceID().String(), spanContext.ParentID().String(), spanContext.SpanID().String()
}

//从 context 获取 trace 信息
func GetTraceInfoFromContext(ctx context.Context) (traceId, pSpanId, spanId string) {
	if ctx == nil {
		return "", "", ""
	}
	//SpanFromContext返回先前与 ctx 相关联的 Span ，如果找不到 Span，则返回 nil。
	//注意：context.Context！= SpanContext：前者是Go的进程内上下文传播机制，后者是OpenTracing的每个 span 标识和 baggage 信息。
	span := opentracing.SpanFromContext(ctx)
	return getTraceInfoFromSpan(span)
}

//jaeger 中间件,用于对每个请求生成一个span，并使用opentracing生成的spanID作为日志的spanID
func JaegerTrace(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request().Header)
		parentSpanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			panic(err)
		}
		operationName := c.Request().RequestURI
		//如果是 get，就得截取 ？ 之前的
		if c.Request().Method == http.MethodGet {
			pos := strings.Index(operationName, "?")
			if pos != -1 {
				operationName = operationName[:pos]
			}
		}
		//创建一个有parentSpan的Span，并设置 span 种类，ip，url等
		span := opentracing.StartSpan(operationName,
			ext.RPCServerOption(parentSpanContext),
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
		//将context关联span
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		//获取 traceId,pSpanId,spanId 后存储到 goroutine 中
		traceId, pSpanId, spanId := GetTraceInfoFromContext(ctx)
		gls.SetGlsWithCtx(traceId, pSpanId, spanId, ctx, func() {
			err = next(c)
		})
		return err
	}
}
