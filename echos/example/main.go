package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/echos"
	"github.com/witwoywhy/go-cores/errs"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/vipers"
)

func init() {
	vipers.Init()
}

func main() {
	app := echos.New[contexts.RouteContext]()

	app.UseMiddleware(app.Log())
	app.UseMiddleware(app.Error())

	app.Register(
		http.MethodGet,
		"/logger",
		app.WithLogger(func(ctx *echo.Context, rctx *contexts.RouteContext, l logger.Logger) {
			l.Info("HANDLE WITH LOGGER")
			ctx.JSON(http.StatusOK, "hello")
		}),
	)

	app.UseMiddleware(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("rctx", &contexts.RouteContext{})
			return next(c)
		}
	})

	app.Register(
		http.MethodPost,
		"/context",
		app.WithRouteContext(func(ctx *echo.Context, rctx *contexts.RouteContext, l logger.Logger) {
			l.Info("HANDLE WITH ROUTE CONTEXT")
			ctx.JSON(http.StatusOK, map[string]string{"hello": "world"})
		}),
	)

	app.Register(
		http.MethodGet,
		"/400",
		func(ctx *echo.Context) error {
			return errs.NewBadRequestError()
		},
	)

	app.ListenAndServe(func() {
		fmt.Println("execute close func")
	})
}
