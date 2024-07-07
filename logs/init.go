package logs

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Init() {
	var config Config
	if err := viper.UnmarshalKey("log", &config); err != nil {
		panic(fmt.Errorf("failed to loaded log config: %v", err))
	}

	var level slog.Level
	switch strings.ToLower(config.Level) {
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

	if config.MaskingList != "" {
		maskingList := strings.Split(config.MaskingList, "|")
		for _, v := range maskingList {
			MaskingList[v] = true
		}
	}
}
