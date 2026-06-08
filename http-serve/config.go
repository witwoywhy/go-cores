package httpserve

import "time"

type HTTPServe struct {
	Port             string   `mapstructure:"port"`
	IgnoreLogBody    []string `mapstructure:"ignore_log_body"`
	ErrorCodeMapping string   `mapstructure:"error_code_mapping"`
	CORS             CORS     `mapstructure:"cors"`
}

type CORS struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}
