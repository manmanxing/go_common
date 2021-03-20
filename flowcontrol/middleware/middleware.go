package middleware

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/labstack/echo"
	"net/http"
)

//参考：https://github.com/alibaba/sentinel-golang/blob/master/pkg/adapters/echo/middleware.go
//todo 针对QPS 设置警告值
func HttpFlowControl(opts ...Option) echo.MiddlewareFunc {
	options := evaluateOptions(opts)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) (err error) {
			resourceName := context.Request().Method + ":" + context.Path()
			if options.resourceExtract != nil {
				resourceName = options.resourceExtract(context)
			}
			entry, blockErr := sentinel.Entry(
				resourceName,
				sentinel.WithResourceType(base.ResTypeWeb),
				sentinel.WithTrafficType(base.Inbound),
			)
			if blockErr != nil {
				if options.blockFallback != nil {
					err = options.blockFallback(context)
				} else {
					//默认的 error 返回
					err = context.JSON(http.StatusTooManyRequests, "请求次数太多,请稍后重试")
				}
				return err
			}

			defer entry.Exit()
			err = next(context)
			return err
		}
	}
}