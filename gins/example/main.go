package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/errs"
	"github.com/witwoywhy/go-cores/gins"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/vipers"
)

func init() {
	vipers.Init()
}

func main() {
	app := gins.New()

	app.UseMiddleware(app.Log())
	app.UseMiddleware(app.Error())

	app.Register(
		http.MethodGet,
		"/logger",
		app.WithParseLogger(func(ctx *gin.Context, l logger.Logger) {
			l.Info("HANDLE WITH LOGGER")
		}),
	)

	app.Register(
		http.MethodGet,
		"/context",
		app.WithParseRouteContext(func(ctx *gin.Context, rctx *contexts.RouteContext, l logger.Logger) {
			l.Info("HANDLE WITH ROUTE CONTEXT")
		}),
	)

	app.Register(
		http.MethodGet,
		"/400",
		app.WithParseLogger(func(ctx *gin.Context, l logger.Logger) {
			ctx.Error(errs.NewBadRequestError())
		}),
	)

	app.ListenAndServe(func() {
		fmt.Println("execute close func")
	})
}
