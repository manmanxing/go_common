package gls

import (
	"context"

	"github.com/jtolds/gls"
)

var (
	mgr         = gls.NewContextManager()
	traceIDKey  = "trace_id"
	pSpanIDKey  = "p_span_id"
	spanIDKey   = "span_id"
	traceCtxKey = "trace_ctx_key"
)

func SetGls(traceID, pSpanID, spanID string, cb func()) {
	mgr.SetValues(gls.Values{traceIDKey: traceID, pSpanIDKey: pSpanID, spanIDKey: spanID}, cb)
}
func SetGlsWithCtx(traceID, pSpanID, spanID string, ctx context.Context, cb func()) {
	mgr.SetValues(gls.Values{traceIDKey: traceID, pSpanIDKey: pSpanID, spanIDKey: spanID, traceCtxKey: ctx}, cb)
}
func GetTraceInfo() (traceID string, pSpanID string, spanID string) {
	trace, ok := mgr.GetValue(traceIDKey)
	if ok {
		traceID = trace.(string)
	}
	pSpan, ok := mgr.GetValue(pSpanIDKey)
	if ok {
		pSpanID = pSpan.(string)
	}
	span, ok := mgr.GetValue(spanIDKey)
	if ok {
		spanID = span.(string)
	}
	return
}
func TraceCtx() (ctx context.Context) {
	traceCtx, ok := mgr.GetValue(traceCtxKey)
	if ok {
		ctx = traceCtx.(context.Context)
	} else {
		ctx = context.Background()
	}
	return
}
