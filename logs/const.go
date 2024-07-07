package logs

import (
	"log/slog"
	"os"
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
	L = New(map[string]any{})
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
