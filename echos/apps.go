package echos

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/errs"
	httpserve "github.com/witwoywhy/go-cores/http-serve"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
)

type App[RouteContext any] interface {
	Register(method string, relativePath string, handlers ...echo.HandlerFunc)
	UseMiddleware(middleware ...echo.MiddlewareFunc)
	WithRouteContext(handle HandleWithRouteContextLogger[RouteContext]) echo.HandlerFunc
	WithLogger(handle HandleWithLogger) echo.HandlerFunc
	ListenAndServe(closeFunc func())

	// middleware
	Log() echo.MiddlewareFunc
	Error() echo.MiddlewareFunc
}

type app[RouteContext any] struct {
	echo *echo.Echo

	config        httpserve.HTTPServe
	ignoreLogBody map[string]bool
	errorMapping  errs.ErrorCodeMapping
}

func New[RouteContext any]() App[RouteContext] {
	var config httpserve.HTTPServe
	if err := viper.UnmarshalKey("http_serve", &config); err != nil {
		panic(fmt.Errorf("failed to loaded [http_serve] config: %v", err))
	}

	a := &app[RouteContext]{
		echo:          echo.New(),
		config:        config,
		ignoreLogBody: map[string]bool{},
		errorMapping:  errs.ParseToErrorCodeMapping(config.ErrorCodeMapping, logs.L),
	}

	for _, v := range config.IgnoreLogBody {
		a.ignoreLogBody[v] = true
	}

	return a
}

func (a *app[RouteContext]) ListenAndServe(closeFunc func()) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Port),
		Handler: a.echo,
	}

	go func() {
		logs.L.Infof("Listen: %s", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.L.Errorf("listen error: %v", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logs.L.Info("start shutdown service ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logs.L.Errorf("shutdown service error: %v", err)
	}

	<-ctx.Done()

	if closeFunc != nil {
		closeFunc()
	}

	logs.L.Info("service shutdown")
}

func (a *app[RouteContext]) Register(method string, relativePath string, handlers ...echo.HandlerFunc) {
	switch method {
	case http.MethodGet:
		a.echo.GET(relativePath, handlers[0])
	case http.MethodPost:
		a.echo.POST(relativePath, handlers[0])
	case http.MethodPatch:
		a.echo.PATCH(relativePath, handlers[0])
	case http.MethodPut:
		a.echo.PUT(relativePath, handlers[0])
	case http.MethodDelete:
		a.echo.DELETE(relativePath, handlers[0])
	}
}

func (a *app[RouteContext]) UseMiddleware(middleware ...echo.MiddlewareFunc) {
	a.echo.Use(middleware...)
}

func (a *app[RouteContext]) WithLogger(handle HandleWithLogger) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		l := NewLogFromCtx(ctx)

		handle(ctx, l)
		return nil
	}
}

const (
	rctx            = "rctx"
	rctxNotFound    = "route context not found"
	rctxInvalidType = "route context invalid type"
)

func (a *app[RouteContext]) WithRouteContext(handle HandleWithRouteContextLogger[RouteContext]) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var l logger.Logger

		ll := ctx.Get("logger")
		if ll == nil {
			l = NewLogFromCtx(ctx)
		} else {
			l = ll.(logger.Logger)
		}

		routeContext := ctx.Get(rctx)
		if routeContext == nil {
			l.Error(rctxNotFound)
			return errs.NewCustom(http.StatusInternalServerError, errs.Err5000500, rctxNotFound, "")
		}

		typedRouteContext, ok := routeContext.(*RouteContext)
		if !ok {
			l.Error(rctxInvalidType)
			return errs.NewCustom(http.StatusInternalServerError, errs.Err5000500, rctxInvalidType, "")
		}

		handle(ctx, typedRouteContext, l)
		return nil
	}
}

func (a *app[RouteContext]) Error() echo.MiddlewareFunc {
	return Error(a.errorMapping)
}

func (a *app[RouteContext]) Log() echo.MiddlewareFunc {
	return Log(a.ignoreLogBody)
}
