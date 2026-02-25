package logasync

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"sync"
	"time"
)

type JsonHandler struct {
	slog.Handler
	l *log.Logger

	ch   chan record
	wg   sync.WaitGroup
	once sync.Once
}

type record struct {
	ctx context.Context
	r   slog.Record
}

func NewJsonHandler(
	out io.Writer,
	options *slog.HandlerOptions,
) *JsonHandler {
	h := &JsonHandler{
		Handler: slog.NewJSONHandler(out, options),
		l:       log.New(out, "", 0),
		ch:      make(chan record, 1024),
		wg:      sync.WaitGroup{},
		once:    sync.Once{},
	}
	h.wg.Add(1)
	go h.worker()
	return h
}

func (h *JsonHandler) worker() {
	defer h.wg.Done()
	for r := range h.ch {
		fields := map[string]any{
			"timestamp": r.r.Time.UnixNano(),
			"datetime":  r.r.Time.Format(time.RFC3339Nano),
			"severity":  r.r.Level,
		}

		r.r.Attrs(func(a slog.Attr) bool {
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
			fields[a.Key] = a.Value.Any()
			return true
		})

		b, err := json.Marshal(fields)
		if err != nil {
			h.l.Println(err)
		}

		h.l.Println(string(b))
	}
}

func (h *JsonHandler) Handle(ctx context.Context, r slog.Record) error {
	clone := r.Clone()
	select {
	case h.ch <- record{ctx: ctx, r: clone}:
	default:
	}
	return nil
}

func (h *JsonHandler) Shutdown(ctx context.Context) error {
	h.once.Do(func() { close(h.ch) })

	done := make(chan struct{})
	go func() { h.wg.Wait(); close(done) }()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
