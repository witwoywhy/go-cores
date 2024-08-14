package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/logger"
)

type HandleWithRouteContextLogger func(ctx *gin.Context, rctx *contexts.RouteContext, l logger.Logger)

type HandleWithLogger func(ctx *gin.Context, l logger.Logger)
