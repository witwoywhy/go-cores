package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
)

func GetIDByKey(key string, ctx *gin.Context) string {
	id := ctx.GetHeader(key)
	if id == "" {
		id = uuid.NewString()
		ctx.Request.Header.Set(key, id)
	}

	return id
}

func NewLogFromCtx(ctx *gin.Context) logger.Logger {
	return logs.New(map[string]any{
		apps.TraceID: GetIDByKey(apps.TraceID, ctx),
		apps.SpanID:  GetIDByKey(apps.SpanID, ctx),
	})
}

// func NewLogWithTraceFromCtx(ctx *gin.Context) logger.Logger {
// 	span, ok := ctx.Get("spanTracer")
// 	if ok {

// 	}
// }
