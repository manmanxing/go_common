package gls

import (
	"context"

	"github.com/jtolds/gls"
)

var (
	newContextManager = gls.NewContextManager()
	traceIDKey        = "trace_id"
	pSpanIDKey        = "p_span_id"
	spanIDKey         = "span_id"
	traceCtxKey       = "trace_ctx_key"
)

func SetGls(traceID, pSpanID, spanID string, cb func()) {
	newContextManager.SetValues(gls.Values{traceIDKey: traceID, pSpanIDKey: pSpanID, spanIDKey: spanID}, cb)
}
func SetGlsWithCtx(traceID, pSpanID, spanID string, ctx context.Context, cb func()) {
	newContextManager.SetValues(gls.Values{traceIDKey: traceID, pSpanIDKey: pSpanID, spanIDKey: spanID, traceCtxKey: ctx}, cb)
}
func GetTraceInfo() (traceID string, pSpanID string, spanID string) {
	trace, ok := newContextManager.GetValue(traceIDKey)
	if ok {
		traceID = trace.(string)
	}
	pSpan, ok := newContextManager.GetValue(pSpanIDKey)
	if ok {
		pSpanID = pSpan.(string)
	}
	span, ok := newContextManager.GetValue(spanIDKey)
	if ok {
		spanID = span.(string)
	}
	return
}
func TraceCtx() (ctx context.Context) {
	traceCtx, ok := newContextManager.GetValue(traceCtxKey)
	if ok {
		ctx = traceCtx.(context.Context)
	} else {
		ctx = context.Background()
	}
	return
}
