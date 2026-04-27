package apps

import (
	"fmt"

	"github.com/spf13/viper"
)

var Config ConfigInfo

type ConfigInfo struct {
	Name      string    `mapstructure:"name"`
	Env       string    `mapstructure:"env"`
	TimeZone  string    `mapstructure:"time_zone"`
	HTTPServe HTTPServe `mapstructure:"http_serve"`
}

type HTTPServe struct {
	Port             string   `mapstructure:"port"`
	IgnoreLogBody    []string `mapstructure:"ignore_log_body"`
	ErrorCodeMapping string   `mapstructure:"error_code_mapping"`
}

func Init() *ConfigInfo {
	if err := viper.UnmarshalKey("app", &Config); err != nil {
		panic(fmt.Errorf("failed to loaded [app] config: %v", err))
	}

	return &Config
}
