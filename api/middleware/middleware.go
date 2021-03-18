package middleware

import (
	"bytes"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"time"
)

//自定义请求流程访问日志
func AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		//这里读取出 body，然后保存起来
		body, _ := ioutil.ReadAll(ctx.Request().Body)
		SetBody(ctx, body)
		//把刚刚读出来的再写进去
		ctx.Request().Body = ioutil.NopCloser(bytes.NewReader(body))
		start := time.Now()
		accessLog(ctx, "access_start", time.Since(start), GetBody(ctx), GetDataIn(ctx), GetDataOut(ctx))
		defer func() {
			accessLog(ctx, "access_end", time.Since(start), GetBody(ctx), GetDataIn(ctx), GetDataOut(ctx))
		}()
		return next(ctx)
	}
}

func accessLog(ctx echo.Context, accessType string, dur time.Duration, body []byte, dataIn, dataOut interface{}) {
	log.Println("type", accessType,
		"ip", ctx.RealIP(), //todo 这里是获取上游还是用户真实ip
		"method", ctx.Request().Method,
		"path", ctx.Request().URL.Path,
		"query", ctx.Request().URL.RawQuery,
		"body", string(body),
		"input", dataIn,
		"output", dataOut,
		"referer", ctx.Request().Referer(),
		"time(ms)", int64(dur/time.Millisecond),
		"request_header", ctx.Request().Header,
	)
}