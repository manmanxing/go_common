package middleware

import (
	"github.com/labstack/echo"
)

func SetDataIn(c echo.Context, dataIn interface{}) {
	c.Set("dataIn", dataIn)
}

func GetDataIn(c echo.Context) (dataIn interface{}) {
	return c.Get("dataIn")
}

func SetDataOut(c echo.Context, dataOut interface{}) {
	c.Set("dataOut", dataOut)
}

func GetDataOut(c echo.Context) (dataOut interface{}) {
	return c.Get("dataOut")
}

func SetBody(c echo.Context, body []byte) {
	c.Set("reqbody", body)
}

func GetBody(c echo.Context) (body []byte) {
	b := c.Get("reqbody")
	if b == nil {
		return nil
	}
	return b.([]byte)
}
