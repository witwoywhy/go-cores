package logs

import (
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/pubsub"
)

var (
	SL *slog.Logger = slog.New(
		NewJsonHandler(
			viper.GetString("app.name"),
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

	Config   ConfigInfo
	producer pubsub.Producer
)

var maskingList = map[string]bool{
	"username":  true,
	"password":  true,
	"email":     true,
	"firstName": true,
	"lastName":  true,
}

const (
	startMessageFmt = "START | %s"
	endMessageFmt   = "END | %s | %v"
)
