package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
	"github.com/witwoywhy/go-cores/tracers"
	"go.opentelemetry.io/otel/attribute"
	ott "go.opentelemetry.io/otel/trace"
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

func NewRequestLog(
	ctx *gin.Context,
	url string,
) (logger.Logger, ott.Span) {
	if logs.LogConfig.IsEnableTracer {
		tcCtx, span := tracers.Trace.Start(ctx, url)

		var (
			traceId = ctx.GetHeader(apps.TraceID)
			spanId  = ctx.GetHeader(apps.SpanID)
		)

		if traceId == "" {
			traceId = span.SpanContext().TraceID().String()
		}

		if spanId == "" {
			spanId = span.SpanContext().SpanID().String()
		}

		span.SetAttributes(attribute.String(apps.TraceID, traceId))
		span.SetAttributes(attribute.String(apps.SpanID, spanId))

		l := &logs.Log{
			Ctx:  tcCtx,
			Span: span,
			Information: map[string]any{
				apps.TraceID: traceId,
				apps.SpanID:  spanId,
			},
		}

		return l, span
	}

	return NewLogFromCtx(ctx), nil
}
