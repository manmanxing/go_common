package config

import (
	"encoding/json"

	"github.com/manmanxing/go_common/beacon/etcd"
)

const (
	//流量控制 etcd 配置路径
	FlowControlEtcdPath = "xxx"
)

//限流配置
type FcConfig struct {
	ServiceName string                  `json:"service_name"`  //服务名
	Rules       map[string]FcConfigInfo `json:"rules"`         //规则，k：request path，v：对应的限流配置
	DefaultRule FcConfigInfo            `json:"default_rule"`  //默认的限流配置
	DingTalkUrl string                  `json:"ding_talk_url"` //钉钉报警 url
}

type FcConfigInfo struct {
	//Resource             string  `json:"resource"`//资源名，即规则的作用目标。
	//MetricType           int32   `json:"metric_type"`//目标指标类型，0 并发，1 QPS
	Count float64 `json:"count"` //阈值
	//RelationStrategy     int32   `json:"relation_strategy"`//关系限流策略，0 使用当前规则的resource， 1 使用关联的resource（在 RefResource 里定义）
	ControlBehavior   int32  `json:"control_behavior"`
	WarmUpPeriodSec   uint32 `json:"warm_up_period_sec"`   //预热的时间长度
	MaxQueueingTimeMs uint32 `json:"max_queueing_time_ms"` //匀速排队的最大等待时间
	//ClusterMode          bool    `json:"cluster_mode"`
	//ClusterThresholdMode uint32  `json:"cluster_threshold_mode"`
}

var HttpControlInfo *FcConfig

/**
{
    "service":"go_center_gateway",
    "rules":{
        "GET:/":{
            "key":"GET:/",
            "threshold_warn":400,
            "threshold_control":4000,
            "control_type":"REJECT"
        },
        "POST:/centerpay/wxpay":{
            "key":"POST:/centerpay/wxpay",
            "threshold_warn":50,
            "threshold_control":500,
            "control_type":"REJECT"
        }
    },
    "default_rule":{
        "threshold_warn":100,
        "threshold_control":2000,
        "control_type":"REJECT"
    },
    "dingtalk_url":"https://oapi.dingtalk.com/robot/send?access_token=dc405ef8d8797e5a0f525549f898f6ed20f8ee03653d652bd031981a70dd7a86"
}
*/

func init() {
	bt, err := etcd.GetValue(FlowControlEtcdPath)
	if err != nil {
		panic(err)
	}
	HttpControlInfo = &FcConfig{
		ServiceName: "",
		Rules:       make(map[string]FcConfigInfo),
		DefaultRule: FcConfigInfo{},
		DingTalkUrl: "",
	}

	err = json.Unmarshal([]byte(bt), HttpControlInfo)
	if err != nil {
		panic(err)
	}
}
