package reqs

import "time"

type Config struct {
	BaseUrl                  string        `mapstructure:"baseUrl"`
	Timeout                  time.Duration `mapstructure:"timeout"`
	EnableInsecureSkipVerify bool          `mapstructure:"enableInsecureSkipVerify"`
	Api                      `mapstructure:",squash"`
}

type Api struct {
	Url    string `mapstructure:"url"`
	Method string `mapstructure:"method"`
}
