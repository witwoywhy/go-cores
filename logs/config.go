package logs

type Config struct {
	Level       string `mapstructure:"level"`
	MaskingList string `mapstructure:"maskingList"`
}
