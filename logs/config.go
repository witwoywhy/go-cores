package logs

type ConfigInfo struct {
	Level       string `mapstructure:"level"`
	MaskingList string `mapstructure:"masking_list"`
	TracerUrl   string `mapstructure:"tracer_url"`

	IsEnableTracer bool
}
