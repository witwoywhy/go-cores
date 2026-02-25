package logs

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/kafka"
)

type KafkaHandler struct {
	slog.Handler
	options *slog.HandlerOptions

	kafka *kafka.Producer
}

func NewKafkaHandler(
	kafka *kafka.Producer,
	options *slog.HandlerOptions,
) *KafkaHandler {
	return &KafkaHandler{
		Handler: slog.NewJSONHandler(nil, options),
		kafka:   kafka,
	}
}

func (h *KafkaHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := map[string]any{
		"timestamp": r.Time.UnixNano(),
		"datetime":  r.Time.Format(time.RFC3339Nano),
		"severity":  r.Level,
	}

	r.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() == slog.KindAny {
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
		}
		fields[a.Key] = a.Value.Any()
		return true
	})

	return h.kafka.Publish(fields[apps.TraceID].(string), fields)
}
