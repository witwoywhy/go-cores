package logs

import (
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
)

var (
	SL *slog.Logger = slog.New(
		NewJsonHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)
	L = New(map[string]any{
		apps.TraceID: uuid.NewString(),
		apps.SpanID:  uuid.NewString(),
	})
)

const (
	Message     = "message"
	MaskingChar = "*"
)

var MaskingList = map[string]bool{
	"username":  true,
	"password":  true,
	"email":     true,
	"firstName": true,
	"lastName":  true,
}
