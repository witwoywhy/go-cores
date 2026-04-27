package tracers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	ot "go.opentelemetry.io/otel/trace"
)

var Tracer *TracerInfo

type TracerInfo struct {
	Provider *trace.TracerProvider
	Trace    ot.Tracer
	Context  context.Context
	Cancel   context.CancelFunc
}

func Init(url string) {
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpointURL(url))
	if err != nil {
		panic(fmt.Errorf("failed when init otel exporter: %v", err))
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(viper.GetString("app.name")),
				attribute.String("environment", viper.GetString("app.env")),
			),
		),
	)

	otel.SetTracerProvider(provider)
	ctx, cancel := context.WithCancel(context.Background())

	Tracer = &TracerInfo{
		Provider: provider,
		Trace:    otel.Tracer(viper.GetString("app.name")),
		Context:  ctx,
		Cancel:   cancel,
	}
}

var onceShutdown sync.Once

func Shutdown() error {
	var err error
	onceShutdown.Do(func() {
		if Tracer != nil && Tracer.Context != nil {
			ctx, cancel := context.WithTimeout(Tracer.Context, 5*time.Second)
			defer cancel()

			err = Tracer.Provider.Shutdown(ctx)

			Tracer.Cancel()
		}
	})
	if err != nil {
		return err
	}

	return nil
}
