package middleware

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo"
	"github.com/manmanxing/go_center_common/api/errorResp"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
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

//健康监测
func HealthCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		eCtx := echoContext{context}
		realIp := eCtx.RealIP()
		//这里可以根据 ip 和 method 筛选
		if context.Request().URL.Path == "/" && strings.HasPrefix(realIp, "100.") {
			switch context.Request().Method {
			case http.MethodGet:
				return context.String(http.StatusOK, "Hello, World!")
			case http.MethodHead:
			default:
				return context.String(http.StatusNotFound, "Not Found")
			}
		}
		return next(context)
	}
}

//异常恢复
//这里使用自定义的统一 error 格式
func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		defer func() {
			if e := recover(); e != nil {
				var recv error
				switch r := e.(type) {
				case error:
					recv = r
				default:
					recv = fmt.Errorf("%v", r)
				}
				//debug.PrintStack()
				stack := make([]byte, 4<<10)         //4 KB
				length := runtime.Stack(stack, true) //打印所有 goroutine
				fmt.Printf("[PANIC RECOVER] %v %s\n", recv, stack[:length])
				//改为统一错误输出
				err = errorresp.ServerError
			}
		}()
		err = next(context)
		return err
	}
}
