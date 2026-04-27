package logs

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/pubsub"
)

type ProducerHandler struct {
	slog.Handler
	serviceName string
	producer    pubsub.Producer
}

func NewProducerHandler(
	serviceName string,
	publisher pubsub.Producer,
	options *slog.HandlerOptions,
) *ProducerHandler {
	return &ProducerHandler{
		serviceName: serviceName,
		Handler:     slog.NewJSONHandler(nil, options),
		producer:    publisher,
	}
}

func (h *ProducerHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := map[string]any{
		"service":   h.serviceName,
		"timestamp": r.Time.UnixNano(),
		"datetime":  r.Time.Format(time.RFC3339Nano),
		"severity":  r.Level,
	}

	r.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() != slog.KindAny {
			fields[a.Key] = a.Value.Any()
			return true
		}

		m, ok := a.Value.Any().(map[string]any)
		if !ok {
			b, err := json.Marshal(a.Value.Any())
			if err != nil {
				return false
			}

			err = json.Unmarshal(b, &m)
			if err != nil {
				return false
			}
		}

		fields[a.Key] = m
		return true
	})

	var message = &pubsub.Message[map[string]any]{
		Context: pubsub.MessageContext{
			TraceID: fields[apps.TraceID].(string),
			SpanID:  fields[apps.SpanID].(string),
		},
		Data: &fields,
	}
	return h.producer.Produce(fields[apps.TraceID].(string), message, L)
}
