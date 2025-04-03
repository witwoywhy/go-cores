package apps

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
	Env      string `mapstructure:"env"`
	TimeZone string `mapstructure:"timeZone"`
}

func InitAppConfig[T any](config *T) {
	if err := viper.UnmarshalKey("app", config); err != nil {
		panic(fmt.Errorf("failed to loaded app config: %v", err))
	}
}
