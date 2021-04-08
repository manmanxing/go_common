package ot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"time"

	qscetcd "code.qschou.com/qschou/go_center_common/beacon/etcd"
	"code.qschou.com/qschou/go_center_common/dlog"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

const (
	JaegerHTTPEndpoint = "http://127.0.0.1:14268/api/traces?format=jaeger.thrift"
)

var jaegerAgentUrl struct {
	JaegerAgentUrl string `json:"jaeger_agent_url"`
}
var jaegerCollectorUrl struct {
	JaegerCollectorUrl string `json:"jaeger_collector_url"`
}
var jaegerProbabilityValue struct {
	JaegerProbabilityValue float64 `json:"jaeger_probability_value"`
}

func init() {
	value, err := qscetcd.Get("root/config/common/opentracing/jaeger_collector")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(value), &jaegerCollectorUrl)
	if err != nil {
		panic(err)
	}

	value, err = qscetcd.Get("root/config/common/opentracing/jaeger_agent")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(value), &jaegerAgentUrl)
	if err != nil {
		panic(err)
	}

	value, err = qscetcd.Get("root/config/common/opentracing/jaeger_probability_value")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(value), &jaegerProbabilityValue)
	if err != nil {
		panic(err)
	}
	if jaegerProbabilityValue.JaegerProbabilityValue < 0.0 || jaegerProbabilityValue.JaegerProbabilityValue > 1.0 {
		panic(fmt.Sprintf("jaegerProbabilityValue must be between 0.0 and 1.0, received %f", jaegerProbabilityValue.JaegerProbabilityValue))
	}
}

type traceLog struct {
}

func (*traceLog) Error(msg string) {
	dlog.Info("opentracing err", msg)
}

func (*traceLog) Infof(msg string, args ...interface{}) {
	dlog.Debug("opentracing info", fmt.Sprintf(msg, args))
}

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
// 可以传入option 自定义采样等，例如：传入config.Sampler(jaeger.NewConstSampler(true)) 或 config.Sampler(jaeger.NewProbabilisticSampler(0.5))
func Init(service string, options ...config.Option) (closer io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: jaegerProbabilityValue.JaegerProbabilityValue,
		},
		Reporter: &config.ReporterConfig{
			QueueSize:          1000,
			LogSpans:           true,
			LocalAgentHostPort: jaegerAgentUrl.JaegerAgentUrl,
			CollectorEndpoint:  jaegerCollectorUrl.JaegerCollectorUrl,
		},
	}
	options = append(options, config.Logger(&traceLog{}), config.Metrics(newMetricRecord()))
	closer, err := cfg.InitGlobalTracer(service, options...)
	if err != nil {
		//panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
		fmt.Println("ERROR: cannot init Jaeger: ", err)
		return
	}
	return
}

func TraceInfoFromSpan(sp opentracing.Span) (traceID string, pSpanID string, spanID string) {
	if sp == nil {
		return
	}
	spctx, ok := sp.Context().(jaeger.SpanContext)
	if !ok {
		return
	}
	traceID = spctx.TraceID().String()
	pSpanID = spctx.ParentID().String()
	spanID = spctx.SpanID().String()
	return
}
func TraceInfo(ctx context.Context) (traceID string, pSpanID string, spanID string) {
	sp := opentracing.SpanFromContext(ctx)
	return TraceInfoFromSpan(sp)
}

type Metric struct {
	CounterRecord   []*CountRecord     `json:"counter_record,omitempty"`
	TimerRecord     []*TimeRecord      `json:"timer_record,omitempty"`
	GaugeRecord     []*GaugeRecord     `json:"gauge_record,omitempty"`
	HistogramRecord []*HistogramRecord `json:"histogram_record,omitempty"`
	NamespaceRecord []*NamespaceRecord `json:"namespace_record,omitempty"`
}

type CountRecord struct {
	Option metrics.Options `json:"option"`
	Count  int64           `json:"count"`
}

func (c *CountRecord) Inc(n int64) {
	c.Count += n
}

type TimeRecord struct {
	Option     metrics.TimerOptions `json:"option"`
	RecordTime time.Duration        `json:"record_time"`
}

func (t *TimeRecord) Record(dur time.Duration) {
	t.RecordTime = dur
}

type GaugeRecord struct {
	Option      metrics.Options `json:"option"`
	RecordValue int64           `json:"record_value"`
}

func (g *GaugeRecord) Update(n int64) {
	g.RecordValue = n
}

type HistogramRecord struct {
	Option      metrics.HistogramOptions `json:"option"`
	RecordValue float64                  `json:"record_value"`
}

func (h *HistogramRecord) Record(n float64) {
	h.RecordValue = n
}

type NamespaceRecord struct {
	Namespace metrics.NSOptions `json:"namespace"`
	SubMetric *Metric           `json:"sub_metric"`
}

func newMetricRecord() *Metric {
	metric := &Metric{}
	metric.recordToLog()
	return metric
}

func (m *Metric) recordToLog() {
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			dlog.Info("log_desc", "opentracing metric record", "metric", m)
		}
	}()
}

func (m *Metric) Counter(metric metrics.Options) metrics.Counter {
	newCounter := &CountRecord{Option: metric}
	m.CounterRecord = append(m.CounterRecord, newCounter)
	return newCounter
}

func (m *Metric) Timer(metric metrics.TimerOptions) metrics.Timer {
	newTimer := &TimeRecord{Option: metric}
	m.TimerRecord = append(m.TimerRecord, newTimer)
	return newTimer
}

func (m *Metric) Gauge(metric metrics.Options) metrics.Gauge {
	newGauge := &GaugeRecord{Option: metric}
	m.GaugeRecord = append(m.GaugeRecord, newGauge)
	return newGauge
}

func (m *Metric) Histogram(metric metrics.HistogramOptions) metrics.Histogram {
	newHistogram := &HistogramRecord{Option: metric}
	m.HistogramRecord = append(m.HistogramRecord, newHistogram)
	return newHistogram
}

func (m *Metric) Namespace(scope metrics.NSOptions) metrics.Factory {
	newNS := &NamespaceRecord{Namespace: scope, SubMetric: &Metric{}}
	m.NamespaceRecord = append(m.NamespaceRecord, newNS)
	return newNS.SubMetric
}
