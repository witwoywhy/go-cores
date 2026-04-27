package logs

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/tracers"
)

func Init(options ...LogConfigOption) {
	var configOption logConfigOption
	for _, option := range options {
		option.apply(&configOption)
	}

	if err := viper.UnmarshalKey("log", &Config); err != nil {
		panic(fmt.Errorf("failed to loaded log config: %v", err))
	}

	var level slog.Level
	switch strings.ToLower(Config.Level) {
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

	if Config.MaskingList != "" {
		for v := range strings.SplitSeq(Config.MaskingList, "|") {
			maskingList[v] = true
		}
	}

	if Config.TracerUrl != "" {
		tracers.Init(Config.TracerUrl)
		NewSpanLog = newSpanLogTracer
		NewSpanLogAction = newSpanLogActionTracer
		Config.IsEnableTracer = true
	}

	producer = configOption.producer
	if producer == nil {
		SL = slog.New(NewJsonHandler(
			viper.GetString("app.name"),
			os.Stdout,
			&slog.HandlerOptions{
				Level: level,
			},
		))
	} else {
		SL = slog.New(NewProducerHandler(
			viper.GetString("app.name"),
			producer,
			&slog.HandlerOptions{
				Level: level,
			},
		))
	}
}

func Shutdown() {
	if producer != nil {
		producer.Shutdown(L)
	}
}
