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
}

func NewJsonHandler(
	out io.Writer,
	options *slog.HandlerOptions,
) *JsonHandler {
	return &JsonHandler{
		Handler: slog.NewJSONHandler(out, options),
		l:       log.New(out, "", 0),
	}
}

func (h *JsonHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := map[string]any{
		"timestamp": r.Time.UnixNano(),
		"datetime":  r.Time.Format(time.RFC3339Nano),
		"level":     r.Level,
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

			masking(m)
			fields[a.Key] = m
			return true
		}
		fields[a.Key] =a.Value.Any()
		return true
	})

	b, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	h.l.Println(string(b))
	return nil
}
