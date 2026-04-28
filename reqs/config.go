package reqs

import "time"

type Config struct {
	BaseUrl                  string        `mapstructure:"base_url"`
	Method                   string        `mapstructure:"method"`
	Url                      string        `mapstructure:"url"`
	Timeout                  time.Duration `mapstructure:"timeout"`
	EnableInsecureSkipVerify bool          `mapstructure:"enable_insecure_skip_verify"`
	EnableIgnoreLogBody      bool          `mapstructure:"enable_ignore_log_body"`
	Addtional                `mapstructure:",squash"`
}

type Addtional struct {
	ApiKey   string `mapstructure:"api_key"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
