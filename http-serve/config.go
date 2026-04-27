package httpserve

type HTTPServe struct {
	Port             string   `mapstructure:"port"`
	IgnoreLogBody    []string `mapstructure:"ignore_log_body"`
	ErrorCodeMapping string   `mapstructure:"error_code_mapping"`
}
