package opentracing

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/manmanxing/go_common/beacon/etcd"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

//jaeger 相关配置信息
var JaegerConfig struct {
	AgentUrl     string  `json:"agent_url"`
	CollectorUrl string  `json:"collector_url"`
	Probability  float64 `json:"probability"`
}

const (
	//jaeger 的etcd配置路径
	jaegerEtcdPath = "xxx"
	//metric 打印日志间隔时间:10分钟
	metricPrintLogTime = 60 * 10
	//在开始删除新span之前可以在内存中保留多少span
	reporterQueueSize = 1000
)

func init() {
	v, err := etcd.GetValue(jaegerEtcdPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(v), &JaegerConfig)
	if err != nil {
		panic(err)
	}
	//强校验jaeger采样率的概率范围
	if JaegerConfig.Probability < 0.0 || JaegerConfig.Probability > 1.0 {
		panic(fmt.Sprintf("jaeger config probability must between 0.0 and 1.0,received %f", JaegerConfig.Probability))
	}
}

//这里可以自定义自己的 tace log
type customizeTaceLog struct {
}

func (*customizeTaceLog) Error(msg string) {
	fmt.Println("customize opentracing err", msg)
}

func (*customizeTaceLog) Infof(msg string, args ...interface{}) {
	fmt.Println("customize opentracing info", fmt.Sprintf(msg, args))
}

//自定义集中式度量系统 metrics
type Metric struct {
	CounterRecord   []*CountRecord     `json:"counter_record,omitempty"`
	TimerRecord     []*TimeRecord      `json:"timer_record,omitempty"`
	GaugeRecord     []*GaugeRecord     `json:"gauge_record,omitempty"`
	HistogramRecord []*HistogramRecord `json:"histogram_record,omitempty"`
	NamespaceRecord []*NamespaceRecord `json:"namespace_record,omitempty"`
}

//自定义实现 metric 的类型
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

//生成 metric 对象
func GetNewMetric() *Metric {
	m := new(Metric)
	m.PrintLog()
	return m
}

//这里可以设置定时打印 metric 信息
func (m *Metric) PrintLog() {
	go func() {
		for {
			time.Sleep(metricPrintLogTime * time.Second)
			fmt.Println("opentracing metric info", m)
		}
	}()
}

//返回一个 Jaeger Trace，可以传入option 自定义采样
func NewTracer(serviceName string, options ...config.Option) (closer io.Closer, err error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: JaegerConfig.Probability,
		},
		Reporter: &config.ReporterConfig{
			QueueSize:          reporterQueueSize,
			LogSpans:           true,
			LocalAgentHostPort: JaegerConfig.AgentUrl,
			CollectorEndpoint:  JaegerConfig.CollectorUrl,
		},
	}
	//这里组装自定义的 log 与 metric
	options = append(options, config.Logger(&customizeTaceLog{}), config.Metrics(GetNewMetric()))
	return cfg.InitGlobalTracer(serviceName, options...)
}
