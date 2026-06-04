package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/logger"
)

type HandleWithRouteContextLogger[T any] func(ctx *gin.Context, rctx *T, l logger.Logger)

type HandleWithLogger func(ctx *gin.Context, rctx *contexts.RouteContext, l logger.Logger)
