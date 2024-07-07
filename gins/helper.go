package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
)

func getIDByKey(key string, ctx *gin.Context) string {
	id := ctx.GetHeader(key)
	if id == "" {
		id = uuid.NewString()
		ctx.Header(key, id)
	}

	return id
}

func NewCoreLogFromCtx(ctx *gin.Context) logger.CoreLogger {
	return logs.NewCoreLog(map[string]any{
		apps.TraceID: getIDByKey(apps.TraceID, ctx),
		apps.SpanID:  getIDByKey(apps.SpanID, ctx),
	})
}
