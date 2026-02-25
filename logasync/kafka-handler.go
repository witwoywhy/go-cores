package logasync

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

	kafka *kafka.Producer
	level slog.Level
}

func NewKafkaHandler(
	kafka *kafka.Producer,
	level slog.Level,
) *KafkaHandler {
	return &KafkaHandler{
		Handler: slog.NewJSONHandler(nil, &slog.HandlerOptions{Level: level}),
		kafka:   kafka,
		level:   level,
	}
}

func (h *KafkaHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
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
