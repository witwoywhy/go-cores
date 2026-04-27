package gins

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/errs"
	httpserve "github.com/witwoywhy/go-cores/http-serve"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
)

type GinApps interface {
	Register(method string, relativePath string, handlers ...gin.HandlerFunc)
	UseMiddleware(middleware ...gin.HandlerFunc)
	WithParseRouteContext(handle HandleWithRouteContextLogger) gin.HandlerFunc
	WithParseLogger(handle HandleWithLogger) gin.HandlerFunc
	ListenAndServe(closeFunc func())

	// middleware
	Log() gin.HandlerFunc
	Error() gin.HandlerFunc
}

type app struct {
	gin *gin.Engine

	config        httpserve.HTTPServe
	ignoreLogBody map[string]bool
	errorMapping  errs.ErrorCodeMapping
}

func New() GinApps {
	var config httpserve.HTTPServe
	if err := viper.UnmarshalKey("http_serve", &config); err != nil {
		panic(fmt.Errorf("failed to loaded [http_serve] config: %v", err))
	}

	a := &app{
		gin:           gin.New(),
		config:        config,
		ignoreLogBody: map[string]bool{},
		errorMapping:  errs.ParseToErrorCodeMapping(config.ErrorCodeMapping, logs.L),
	}

	for _, v := range config.IgnoreLogBody {
		a.ignoreLogBody[v] = true
	}

	return a
}

func (a *app) ListenAndServe(closeFunc func()) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Port),
		Handler: a.gin,
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

func (a *app) Register(method string, relativePath string, handlers ...gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		a.gin.GET(relativePath, handlers...)
	case http.MethodPost:
		a.gin.POST(relativePath, handlers...)
	case http.MethodPatch:
		a.gin.PATCH(relativePath, handlers...)
	case http.MethodPut:
		a.gin.PUT(relativePath, handlers...)
	case http.MethodDelete:
		a.gin.DELETE(relativePath, handlers...)
	}
}

func (a *app) UseMiddleware(middleware ...gin.HandlerFunc) {
	a.gin.Use(middleware...)
}

func (a *app) WithParseLogger(handle HandleWithLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := NewLogFromCtx(ctx)

		handle(ctx, l)
	}
}

const (
	rctx         = "rctx"
	rctxNotFound = "route context not found"
)

func (a *app) WithParseRouteContext(handle HandleWithRouteContextLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var l logger.Logger

		ll, ok := ctx.Get("logger")
		if !ok {
			l = NewLogFromCtx(ctx)
		} else {
			l = ll.(logger.Logger)
		}

		routeContext, ok := ctx.Get(rctx)
		if !ok {
			l.Error("rctx not found")
			ctx.Error(errs.NewCustom(http.StatusInternalServerError, errs.Err50000, rctxNotFound, ""))
			ctx.Abort()
		}

		handle(ctx, routeContext.(*contexts.RouteContext), l)
	}
}
