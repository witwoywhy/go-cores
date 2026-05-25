package echos

import (
	"github.com/labstack/echo/v5"
	"github.com/witwoywhy/go-cores/logger"
)

type HandleWithRouteContextLogger[T any] func(ctx *echo.Context, rctx *T, l logger.Logger)

type HandleWithLogger func(ctx *echo.Context, l logger.Logger)
