package echos

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/enum/language"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
	"github.com/witwoywhy/go-cores/tracers"
	"go.opentelemetry.io/otel/attribute"
	ot "go.opentelemetry.io/otel/trace"
)

func getIDByKey(key string, ctx *echo.Context) string {
	id := ctx.Request().Header.Get(key)
	if id == "" {
		id = uuid.NewString()
		ctx.Request().Header.Set(key, id)
	}

	return id
}

func NewLogFromCtx(ctx *echo.Context) logger.Logger {
	return logs.New(map[string]any{
		apps.TraceID: getIDByKey(apps.TraceID, ctx),
		apps.SpanID:  getIDByKey(apps.SpanID, ctx),
	})
}

func NewRequestLog(
	ctx *echo.Context,
	url string,
) (logger.Logger, ot.Span) {
	if !logs.Config.IsEnableTracer {
		return NewLogFromCtx(ctx), nil
	}

	tcCtx, span := tracers.Tracer.Trace.Start(ctx.Request().Context(), url)

	var (
		traceID = ctx.Request().Header.Get(apps.TraceID)
		spanID  = ctx.Request().Header.Get(apps.SpanID)
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

func GetLanguage(ctx *echo.Context) language.Language {
	headerLang := ctx.Request().Header.Get("X-Language")
	if headerLang != "" {
		return language.Language(headerLang)
	}

	lang := ctx.Get(apps.Language)
	if lang != nil {
		return lang.(language.Language)
	}

	return language.TH
}
