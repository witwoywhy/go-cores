package logs

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/tracers"
	"go.opentelemetry.io/otel/attribute"
	ot "go.opentelemetry.io/otel/trace"
)

type Log struct {
	Information map[string]any
	args        []any

	ctx  context.Context
	span ot.Span
}

func New(info map[string]any) logger.Logger {
	log := &Log{
		Information: info,
		args:        []any{},
	}

	if log.Information == nil {
		log.Information = map[string]any{}
	}

	for k, v := range log.Information {
		log.args = append(log.args, k, v)
	}

	return log
}

var NewSpanLog func(l logger.Logger) logger.Logger = newSpanLogStd

func newSpanLogStd(l logger.Logger) logger.Logger {
	ll := l.(*Log)

	newLog := &Log{
		Information: maps.Clone(ll.Information),
	}
	newLog.Information[apps.SpanID] = uuid.NewString()
	newLog.AddArgs()

	return newLog
}

func newSpanLogTracer(l logger.Logger) logger.Logger {
	ll := l.(*Log)

	ctx, span := ll.span.TracerProvider().Tracer("").Start(ll.ctx, "")
	span.SetAttributes(attribute.String(apps.TraceID, span.SpanContext().TraceID().String()))
	span.SetAttributes(attribute.String(apps.SpanID, span.SpanContext().SpanID().String()))

	newLog := &Log{
		Information: maps.Clone(ll.Information),
		ctx:         ctx,
		span:        span,
	}
	newLog.Information[apps.SpanID] = span.SpanContext().SpanID().String()
	newLog.AddArgs()
	return newLog
}

var NewSpanLogAction func(l logger.Logger, action string) (logger.Logger, func()) = newSpanLogActionStd

func newSpanLogActionStd(l logger.Logger, action string) (logger.Logger, func()) {
	now := time.Now()
	ll := newSpanLogStd(l)
	ll.Infof(startMessageFmt, action)
	return ll, func() {
		ll.Infof(endMessageFmt, action, time.Since(now))
	}
}

func newSpanLogActionTracer(l logger.Logger, action string) (logger.Logger, func()) {
	ll := l.(*Log)
	now := time.Now()

	ctx, span := ll.span.TracerProvider().Tracer("").Start(ll.ctx, action)
	span.SetAttributes(attribute.String(apps.TraceID, span.SpanContext().TraceID().String()))
	span.SetAttributes(attribute.String(apps.SpanID, span.SpanContext().SpanID().String()))

	newLog := &Log{
		Information: maps.Clone(ll.Information),
		ctx:         ctx,
		span:        span,
	}
	newLog.Information[apps.SpanID] = span.SpanContext().SpanID().String()
	newLog.AddArgs()

	l = newLog
	l.Infof(startMessageFmt, action)
	return l, func() {
		since := time.Since(now)
		span.SetAttributes(attribute.String(apps.ProcessTime, since.String()))
		span.End()
		l.Infof(endMessageFmt, action, since)
	}
}

func NewLogTracer(
	ctx context.Context,
	action string,
) (logger.Logger, ot.Span) {
	ctx, span := tracers.Tracer.Trace.Start(ctx, action)
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

func NewLog() (logger.Logger, ot.Span) {
	if !Config.IsEnableTracer {
		return New(map[string]any{
			apps.TraceID: uuid.NewString(),
			apps.SpanID:  uuid.NewString(),
		}), nil
	}

	ctx, span := tracers.Tracer.Trace.Start(context.Background(), viper.GetString("app.name"))
	var (
		traceID = span.SpanContext().TraceID().String()
		spanID  = span.SpanContext().SpanID().String()
	)

	span.SetAttributes(attribute.String(apps.TraceID, traceID))
	span.SetAttributes(attribute.String(apps.SpanID, spanID))

	l := &Log{
		Information: map[string]any{
			apps.TraceID: traceID,
			apps.SpanID:  spanID,
		},
	}
	l.AddArgs()
	l.AddTracer(ctx, span)

	return l, span
}

func (l *Log) Debug(obj any) {
	args := slices.Clone(l.args)
	args = append(args, apps.Message, obj)
	SL.Debug("", args...)
}

func (l *Log) Debugf(format string, obj ...any) {
	args := slices.Clone(l.args)
	args = append(args, apps.Message, fmt.Sprintf(format, obj...))
	SL.Debug("", args...)
}

func (l *Log) Error(err any) {
	args := slices.Clone(l.args)
	switch v := err.(type) {
	case error:
		args = append(args, apps.Message, v.Error())
	default:
		args = append(args, apps.Message, v)
	}
	SL.Error("", args...)
}

func (l *Log) Errorf(format string, obj ...any) {
	args := slices.Clone(l.args)
	args = append(args, apps.Message, fmt.Sprintf(format, obj...))
	SL.Error("", args...)
}

func (l *Log) Info(obj any) {
	args := slices.Clone(l.args)
	args = append(args, apps.Message, obj)
	SL.Info("", args...)
}

func (l *Log) Infof(format string, obj ...any) {
	args := slices.Clone(l.args)
	args = append(args, apps.Message, fmt.Sprintf(format, obj...))
	SL.Info("", args...)
}

func (l *Log) Warn(obj any) {
	args := slices.Clone(l.args)
	args = append(args, apps.Message, obj)
	SL.Warn("", args...)
}

func (l *Log) Warnf(format string, obj ...any) {
	args := slices.Clone(l.args)
	args = append(args, apps.Message, fmt.Sprintf(format, obj...))
	SL.Warn("", args...)
}

func (l *Log) JSON(m map[string]any) {
	args := slices.Clone(l.args)
	for k, v := range m {
		args = append(args, k, v)
	}
	SL.Info("", args...)
}

func (l *Log) AddTracer(ctx context.Context, span ot.Span) {
	l.ctx = ctx
	l.span = span
}

func (l *Log) AddArgs() {
	for k, v := range l.Information {
		l.args = append(l.args, k, v)
	}
}
