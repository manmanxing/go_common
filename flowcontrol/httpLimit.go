package flowcontrol

import (
	"fmt"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/labstack/echo"
	"github.com/manmanxing/go_common/flowcontrol/config"
	"github.com/manmanxing/go_common/flowcontrol/middleware"
)

//http接口限流
func InitHttpFlowRate(serviceName string, e *echo.Echo) {
	routes := e.Routes()
	err := sentinel.InitDefault()
	if err != nil {
		fmt.Printf("Unexpected error: %+v \n", err)
		return
	}
	var rules []*flow.FlowRule
	for _, v := range routes {
		resourceName := fmt.Sprintf("%s:%s", v.Method, v.Path)
		//先去查询该 resource
		//若没有找到，就找 path
		//如果还没找到就去使用默认限流配置
		r, ok := config.HttpControlInfo.Rules[resourceName]
		if !ok {
			r, ok = config.HttpControlInfo.Rules[v.Path]
			if !ok {
				r = config.HttpControlInfo.DefaultRule
			}
		}

		rule := &flow.FlowRule{
			Resource:          fmt.Sprintf("%s:%s", v.Method, v.Path),
			MetricType:        flow.QPS,
			Count:             r.Count,
			RelationStrategy:  flow.Direct,
			ControlBehavior:   flow.ControlBehavior(r.ControlBehavior),
			WarmUpPeriodSec:   r.WarmUpPeriodSec,
			MaxQueueingTimeMs: r.MaxQueueingTimeMs,
		}

		rules = append(rules, rule)

	}

	_, err = flow.LoadRules(rules)
	if err != nil {
		return
	}

	e.Use(middleware.HttpFlowControl())
}
