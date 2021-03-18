package util

import (
	"github.com/labstack/echo"
	"net/http/pprof"
	"strings"
)

/*
//main 方法调用
	go func() {
		PPROFListenPort := "32003"
		e := echo.New()
		e.GET("/ping", func(c echo.Context) error {
			return c.String(200, "pong")
		})
		api.Wrap(e)
		err := e.Start(":" + PPROFListenPort)
		if err != nil {
			panic(err)
		}
		dlog.Info("server start", topic, "addr", PPROFListenPort, "version", "v1.0.0")
	}()
*/

func Wrap(e *echo.Echo) {
	WrapGroup("", e.Group(""))
}

func WrapGroup(prefix string, g *echo.Group) {
	routers := []struct {
		Method  string
		Path    string
		Handler echo.HandlerFunc
	}{
		{"GET", "/debug/pprof/", IndexHandler()},
		{"GET", "/debug/pprof/heap", HeapHandler()},
		{"GET", "/debug/pprof/goroutine", GoroutineHandler()},
		{"GET", "/debug/pprof/block", BlockHandler()},
		{"GET", "/debug/pprof/threadcreate", ThreadCreateHandler()},
		{"GET", "/debug/pprof/cmdline", CmdlineHandler()},
		{"GET", "/debug/pprof/profile", ProfileHandler()},
		{"GET", "/debug/pprof/symbol", SymbolHandler()},
		{"POST", "/debug/pprof/symbol", SymbolHandler()},
		{"GET", "/debug/pprof/trace", TraceHandler()},
		{"GET", "/debug/pprof/mutex", MutexHandler()},
	}

	for _, r := range routers {
		switch r.Method {
		case "GET":
			g.GET(strings.TrimPrefix(r.Path, prefix), r.Handler)
		case "POST":
			g.POST(strings.TrimPrefix(r.Path, prefix), r.Handler)
		}
	}
}

func IndexHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Index(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func HeapHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

func GoroutineHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func BlockHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("block").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func ThreadCreateHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func CmdlineHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Cmdline(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func ProfileHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Profile(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func SymbolHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Symbol(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func TraceHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Trace(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

func MutexHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("mutex").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}
