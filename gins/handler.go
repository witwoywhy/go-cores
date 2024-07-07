package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/logger"
)

type HandlerFunc interface {
	HandleFunc | HandlerWithLogFunc
}

type HandleFunc func(ctx *gin.Context)

type HandlerWithLogFunc func(ctx *gin.Context, l logger.Logger)
