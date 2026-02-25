package logasync

var LogConfig Config

type Config struct {
	Level       string  `mapstructure:"level"`
	MaskingList string  `mapstructure:"maskingList"`
	TracerUrl   *string `mapstructure:"tracerUrl"`
	IsPubsub    bool    `mapstructure:"isPubsub"`

	IsEnableTracer bool
}
