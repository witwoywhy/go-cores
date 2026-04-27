package kafka

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
	"github.com/witwoywhy/go-cores/tracers"
	"go.opentelemetry.io/otel/attribute"
	ot "go.opentelemetry.io/otel/trace"
)

func newLogTracer(topic, group, traceID string) (logger.Logger, ot.Span) {
	ctx, span := tracers.Tracer.Trace.Start(context.Background(), fmt.Sprintf("%s:%s", topic, group))

	if traceID == "" {
		traceID = span.SpanContext().SpanID().String()
	}

	var spanID = span.SpanContext().SpanID().String()
	span.SetAttributes(attribute.String(apps.TraceID, traceID))
	span.SetAttributes(attribute.String(apps.SpanID, spanID))

	l := &logs.Log{
		Information: map[string]any{
			apps.TraceID: traceID,
			apps.SpanID:  spanID,
		},
	}
	l.AddArgs()
	l.AddTracer(ctx, span)

	return l, span
}

func newSTDLog(traceID string) (logger.Logger, ot.Span) {
	if traceID == "" {
		traceID = uuid.NewString()
	}

	return logs.New(map[string]any{
		apps.TraceID: traceID,
		apps.SpanID:  uuid.NewString(),
	}), nil
}

func NewLog(topic, group, traceID string) (logger.Logger, ot.Span) {
	if logs.Config.IsEnableTracer {
		return newLogTracer(topic, group, traceID)
	}

	return newSTDLog(traceID)
}
