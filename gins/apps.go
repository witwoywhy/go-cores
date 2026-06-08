package gins

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/errs"
	httpserve "github.com/witwoywhy/go-cores/http-serve"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
)

type App[RouteContext any] interface {
	Register(method string, relativePath string, handlers ...gin.HandlerFunc)
	Group(prefix string, handlers ...gin.HandlerFunc) Group
	UseMiddleware(middleware ...gin.HandlerFunc)
	WithRouteContext(handle HandleWithRouteContextLogger[RouteContext]) gin.HandlerFunc
	WithLogger(handle HandleWithLogger) gin.HandlerFunc
	ListenAndServe(closeFunc func())

	// middleware
	Log() gin.HandlerFunc
	Error() gin.HandlerFunc
}

type Group interface {
	Register(method string, relativePath string, handlers ...gin.HandlerFunc)
	Group(prefix string, handlers ...gin.HandlerFunc) Group
	UseMiddleware(middleware ...gin.HandlerFunc)
}

type app[RouteContext any] struct {
	gin *gin.Engine

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
		gin:           gin.New(),
		config:        config,
		ignoreLogBody: map[string]bool{},
		errorMapping:  errs.ParseToErrorCodeMapping(config.ErrorCodeMapping, logs.L),
	}

	a.UseMiddleware(cors.New(cors.Config{
		AllowOrigins:     a.config.CORS.AllowOrigins,
		AllowMethods:     a.config.CORS.AllowMethods,
		AllowHeaders:     a.config.CORS.AllowHeaders,
		AllowCredentials: a.config.CORS.AllowCredentials,
		MaxAge:           a.config.CORS.MaxAge,
		AllowAllOrigins:  false,
	}))

	for _, v := range config.IgnoreLogBody {
		a.ignoreLogBody[v] = true
	}

	return a
}

func (a *app[RouteContext]) ListenAndServe(closeFunc func()) {
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

func (a *app[RouteContext]) Register(method string, relativePath string, handlers ...gin.HandlerFunc) {
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

func (a *app[RouteContext]) Group(prefix string, handlers ...gin.HandlerFunc) Group {
	return &group{group: a.gin.Group(prefix, handlers...)}
}

func (a *app[RouteContext]) UseMiddleware(middleware ...gin.HandlerFunc) {
	a.gin.Use(middleware...)
}

func (a *app[RouteContext]) WithLogger(handle HandleWithLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := NewLogFromCtx(ctx)

		rctx := &contexts.RouteContext{Ctx: ctx}
		if err := ctx.BindHeader(&rctx); err != nil {
			l.Errorf("failed when bind header: %v", err)
			ctx.Error(errs.NewBadRequestError())
			ctx.Abort()
			return
		}

		handle(ctx, rctx, l)
	}
}

const (
	rctx            = "rctx"
	rctxNotFound    = "route context not found"
	rctxInvalidType = "route context invalid type"
)

func (a *app[RouteContext]) WithRouteContext(handle HandleWithRouteContextLogger[RouteContext]) gin.HandlerFunc {
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
			ctx.Error(errs.NewCustom(http.StatusInternalServerError, errs.Err5000500, rctxNotFound, ""))
			ctx.Abort()
			return
		}

		typedRouteContext, ok := routeContext.(*RouteContext)
		if !ok {
			l.Error("rctx invalid type")
			ctx.Error(errs.NewCustom(http.StatusInternalServerError, errs.Err5000500, rctxInvalidType, ""))
			ctx.Abort()
			return
		}

		handle(ctx, typedRouteContext, l)
	}
}

func (a *app[RouteContext]) Error() gin.HandlerFunc {
	return Error(a.errorMapping)
}

func (a *app[RouteContext]) Log() gin.HandlerFunc {
	return Log(a.ignoreLogBody)
}

type group struct {
	group *gin.RouterGroup
}

func (g *group) Register(method string, relativePath string, handlers ...gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		g.group.GET(relativePath, handlers...)
	case http.MethodPost:
		g.group.POST(relativePath, handlers...)
	case http.MethodPatch:
		g.group.PATCH(relativePath, handlers...)
	case http.MethodPut:
		g.group.PUT(relativePath, handlers...)
	case http.MethodDelete:
		g.group.DELETE(relativePath, handlers...)
	}
}

func (g *group) Group(prefix string, handlers ...gin.HandlerFunc) Group {
	return &group{group: g.group.Group(prefix, handlers...)}
}

func (g *group) UseMiddleware(middleware ...gin.HandlerFunc) {
	g.group.Use(middleware...)
}
