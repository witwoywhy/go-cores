package tracers

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	ot "go.opentelemetry.io/otel/trace"
)

var (
	TraceProvider *trace.TracerProvider
	TraceContext  context.Context
	TraceCancel   context.CancelFunc
	Trace         ot.Tracer
)

func InitTracer(url string) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(fmt.Errorf("failed when init otel exporter: %v", err))
	}

	TraceProvider = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(viper.GetString("app.name")),
				attribute.String("environment", viper.GetString("app.env")),
			),
		),
	)

	TraceContext, TraceCancel = context.WithCancel(context.Background())

	otel.SetTracerProvider(TraceProvider)
	Trace = otel.Tracer(viper.GetString("app.name"))
}

func ShutDown() {
	if TraceContext != nil {
		ctx, cancel := context.WithTimeout(TraceContext, 5*time.Second)
		defer cancel()

		if err := TraceProvider.Shutdown(ctx); err != nil {
			panic(fmt.Errorf("failed when init shutdown tracer: %v", err))
		}

		TraceCancel()
	}
}
