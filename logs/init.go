package logs

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/tracers"
	"github.com/witwoywhy/go-cores/utils"
)

func Init() {
	if err := viper.UnmarshalKey("log", &LogConfig); err != nil {
		panic(fmt.Errorf("failed to loaded log config: %v", err))
	}

	var level slog.Level
	switch strings.ToLower(LogConfig.Level) {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	SL = slog.New(
		NewJsonHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: level,
			},
		),
	)

	if LogConfig.MaskingList != "" {
		maskingList := strings.Split(LogConfig.MaskingList, "|")
		for _, v := range maskingList {
			MaskingList[v] = true
		}
	}

	if LogConfig.TracerUrl != nil {
		tracers.InitTracer(utils.NotNil(LogConfig.TracerUrl))
		NewSpanLogAction = NewSpanLogActionTracer
		LogConfig.IsEnableTracer = true
	}
}
