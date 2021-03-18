package middleware

import (
	"github.com/labstack/echo"
	"net"
	"strings"
)

//针对 context 进行处理
func HookCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		ectx := &echoContext{context}
		return next(ectx)
	}
}

type echoContext struct {
	echo.Context
}

//这里在context 保存请求进出时刻的信息，方便打印请求流程日志
func (e *echoContext) Bind(i interface{}) error {
	SetDataIn(e.Context, i)
	return e.Context.Bind(i)
}

func (e *echoContext) JSON(code int, i interface{}) error {
	SetDataOut(e.Context, i)
	return e.Context.JSON(code, i)
}
func (e *echoContext) String(code int, s string) error {
	SetDataOut(e.Context, s)
	return e.Context.String(code, s)
}
func (e *echoContext) HTML(code int, html string) error {
	SetDataOut(e.Context, html)
	return e.Context.HTML(code, html)
}
func (e *echoContext) XML(code int, i interface{}) error {
	SetDataOut(e.Context, i)
	return e.Context.XML(code, i)
}

//自定义获取真实ip
func (e *echoContext) RealIP() string {
	var ip string
	ip = e.Request().Header.Get(echo.HeaderXRealIP)
	if len(strings.TrimSpace(ip)) > 0 {
		return ip
	}
	ip = e.Request().Header.Get(echo.HeaderXForwardedFor)
	if len(strings.TrimSpace(ip)) > 0 {
		//取第一个
		ipList := strings.Split(ip, ",")
		ip = ipList[0]
		return ip
	}
	ip, _, _ = net.SplitHostPort(e.Request().RemoteAddr)
	return ip
}