package logs

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/witwoywhy/go-cores/logger"
)

type jsonHandler struct {
	slog.Handler
	writer logger.Writer

	addServiceNameFunc func(fields map[string]any)
}

func NewJsonHandler(
	serviceName string,
	writer logger.Writer,
	options *slog.HandlerOptions,

) *jsonHandler {
	handler := &jsonHandler{
		Handler: slog.NewJSONHandler(os.Stdout, options),
		writer:  writer,
	}

	if serviceName == "" {
		handler.addServiceNameFunc = func(fields map[string]any) {}
	} else {
		handler.addServiceNameFunc = func(fields map[string]any) {
			fields["service"] = serviceName
		}
	}

	return handler
}

func (h *jsonHandler) addServiceName(fields map[string]any) {
	h.addServiceNameFunc(fields)
}

func (h *jsonHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := make(map[string]any, 3+r.NumAttrs())
	fields["timestamp"] = r.Time.UnixNano()
	fields["datetime"] = r.Time.Format(time.RFC3339Nano)
	fields["severity"] = r.Level
	h.addServiceName(fields)

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

		masking(m)
		fields[a.Key] = m
		return true
	})

	h.writer.Write(fields, L)
	return nil
}
