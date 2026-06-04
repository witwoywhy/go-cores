package reqs

import "time"

type Config struct {
	BaseUrl                  string            `mapstructure:"base_url"`
	Method                   string            `mapstructure:"method"`
	Url                      string            `mapstructure:"url"`
	Timeout                  time.Duration     `mapstructure:"timeout"`
	EnableInsecureSkipVerify bool              `mapstructure:"enable_insecure_skip_verify"`
	EnableIgnoreLogBody      bool              `mapstructure:"enable_ignore_log_body"`
	Addtional                map[string]string `mapstructure:",remain"`
}
