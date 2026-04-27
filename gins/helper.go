package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/enum/language"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
	"github.com/witwoywhy/go-cores/tracers"
	"go.opentelemetry.io/otel/attribute"
	ot "go.opentelemetry.io/otel/trace"
)

func getIDByKey(key string, ctx *gin.Context) string {
	id := ctx.GetHeader(key)
	if id == "" {
		id = uuid.NewString()
		ctx.Request.Header.Set(key, id)
	}

	return id
}

func NewLogFromCtx(ctx *gin.Context) logger.Logger {
	return logs.New(map[string]any{
		apps.TraceID: getIDByKey(apps.TraceID, ctx),
		apps.SpanID:  getIDByKey(apps.SpanID, ctx),
	})
}

func NewRequestLog(
	ctx *gin.Context,
	url string,
) (logger.Logger, ot.Span) {
	if !logs.Config.IsEnableTracer {
		return NewLogFromCtx(ctx), nil
	}

	tcCtx, span := tracers.Tracer.Trace.Start(ctx, url)

	var (
		traceID = ctx.GetHeader(apps.TraceID)
		spanID  = ctx.GetHeader(apps.SpanID)
	)

	if traceID == "" {
		traceID = span.SpanContext().TraceID().String()
	}

	if spanID == "" {
		spanID = span.SpanContext().SpanID().String()
	}

	span.SetAttributes(attribute.String(apps.TraceID, traceID))
	span.SetAttributes(attribute.String(apps.SpanID, spanID))
	ctx.Set(apps.TraceID, traceID)
	ctx.Set(apps.SpanID, spanID)

	l := &logs.Log{
		Information: map[string]any{
			apps.TraceID: traceID,
			apps.SpanID:  spanID,
		},
	}
	l.AddArgs()
	l.AddTracer(tcCtx, span)

	return l, span
}

func GetLanguage(ctx *gin.Context) language.Language {
	headerLang := ctx.GetHeader("X-Language")
	if headerLang != "" {
		return language.Language(headerLang)
	}

	lang, ok := ctx.Get(apps.Language)
	if ok {
		return lang.(language.Language)
	}

	return language.TH
}
