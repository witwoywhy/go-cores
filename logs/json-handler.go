package logs

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"time"
)

type JsonHandler struct {
	slog.Handler
	l *log.Logger

	addServiceNameFunc func(fields map[string]any)
}

func NewJsonHandler(
	serviceName string,
	out io.Writer,
	options *slog.HandlerOptions,

) *JsonHandler {
	handler := &JsonHandler{
		Handler: slog.NewJSONHandler(out, options),
		l:       log.New(out, "", 0),
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

func (h *JsonHandler) Handle(ctx context.Context, r slog.Record) error {
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

	b, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	h.l.Println(string(b))
	return nil
}

func (h *JsonHandler) addServiceName(fields map[string]any) {
	h.addServiceNameFunc(fields)
}
