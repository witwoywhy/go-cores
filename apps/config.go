package apps

import (
	"fmt"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

func InitAppConfig() {
	var config Config
	if err := viper.UnmarshalKey("app", &config); err != nil {
		panic(fmt.Errorf("failed to loaded app config: %v", err))
	}

	AppConfig = config
}
