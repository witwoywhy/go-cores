package logs

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"go.opentelemetry.io/otel/attribute"
	ot "go.opentelemetry.io/otel/trace"
)

type Log struct {
	Information map[string]any

	ctx  context.Context
	span ot.Span
}

func New(info map[string]any) logger.Logger {
	return &Log{
		Information: info,
	}
}

func NewSpanLog(l logger.Logger) logger.Logger {
	ll := l.(*Log)
	newLog := Log{Information: maps.Clone(ll.Information)}
	newLog.Information[apps.SpanID] = uuid.NewString()
	return &newLog
}

func NewSpanLogAction(l logger.Logger, action string) (logger.Logger, func()) {
	now := time.Now()
	ll := NewSpanLog(l)
	ll.Infof("START | %s", action)
	return ll, func() {
		ll.Infof("END | %s | %v", action, time.Since(now))
	}
}

func NewTracer(
	ctx context.Context,
	trace ot.Tracer,
	action string,
) (logger.Logger, ot.Span) {
	ctx, span := trace.Start(ctx, action)
	span.SetAttributes(attribute.String(apps.TraceID, span.SpanContext().TraceID().String()))
	span.SetAttributes(attribute.String(apps.SpanID, span.SpanContext().SpanID().String()))

	l := &Log{
		ctx:  ctx,
		span: span,
		Information: map[string]any{
			apps.TraceID: span.SpanContext().TraceID().String(),
			apps.SpanID:  span.SpanContext().SpanID().String(),
		},
	}

	return l, span
}

func NewSpanTracer(l logger.Logger, action string) (logger.Logger, func()) {
	ll := l.(*Log)
	now := time.Now()

	ctx, span := ll.span.TracerProvider().Tracer("").Start(ll.ctx, action)
	span.SetAttributes(attribute.String(apps.TraceID, span.SpanContext().TraceID().String()))
	span.SetAttributes(attribute.String(apps.SpanID, span.SpanContext().SpanID().String()))

	newLog := &Log{Information: maps.Clone(ll.Information), ctx: ctx, span: span}
	newLog.Information[apps.SpanID] = span.SpanContext().SpanID().String()

	l = newLog
	l.Infof("START | %s", action)
	return l, func() {
		since := time.Since(now)
		span.SetAttributes(attribute.String("since", since.String()))
		span.End()
		l.Infof("END | %s | %v", action, since)
	}
}

func (l *Log) Debug(obj any) {
	args := infoToArgs(l.Information)
	args = append(args, Message, obj)
	SL.Debug("", args...)
}

func (l *Log) Debugf(format string, obj ...any) {
	args := infoToArgs(l.Information)
	args = append(args, Message, fmt.Sprintf(format, obj...))
	SL.Debug("", args...)
}

func (l *Log) Error(err any) {
	args := infoToArgs(l.Information)
	switch v := err.(type) {
	case error:
		args = append(args, Message, v.Error())
	default:
		args = append(args, Message, v)
	}

	SL.Error("", args...)
}

func (l *Log) Errorf(format string, obj ...any) {
	args := infoToArgs(l.Information)
	args = append(args, Message, fmt.Sprintf(format, obj...))
	SL.Error("", args...)
}

func (l *Log) Info(obj any) {
	args := infoToArgs(l.Information)
	args = append(args, Message, obj)
	SL.Info("", args...)
}

func (l *Log) Infof(format string, obj ...any) {
	args := infoToArgs(l.Information)
	args = append(args, Message, fmt.Sprintf(format, obj...))
	SL.Info("", args...)
}

func (l *Log) Warn(obj any) {
	args := infoToArgs(l.Information)
	args = append(args, Message, obj)
	SL.Warn("", args...)
}

func (l *Log) Warnf(format string, obj ...any) {
	args := infoToArgs(l.Information)
	args = append(args, Message, fmt.Sprintf(format, obj...))
	SL.Warn("", args...)
}

func (l *Log) JSON(m map[string]any) {
	args := infoToArgs(l.Information)
	for k, v := range m {
		args = append(args, k, v)
	}
	SL.Info("", args...)
}

func (l *Log) AddInformation(m map[string]any) {
	for k, v := range m {
		l.Information[k] = v
	}
}
