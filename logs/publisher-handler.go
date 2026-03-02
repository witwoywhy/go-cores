package logs

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/witwoywhy/go-cores/apps"

	"github.com/witwoywhy/go-cores/pubsub"
)

type PublisherHandler struct {
	slog.Handler
	publisher pubsub.Publisher
}

func NewPublisherHandler(
	publisher pubsub.Publisher,
	options *slog.HandlerOptions,
) *PublisherHandler {
	return &PublisherHandler{
		Handler:   slog.NewJSONHandler(nil, options),
		publisher: publisher,
	}
}

func (h *PublisherHandler) Handle(ctx context.Context, r slog.Record) error {
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

	return h.publisher.Publish(fields[apps.TraceID].(string), fields, L)
}
