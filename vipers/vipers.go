package vipers

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	FileNameConfig      = "config"
	FileExtensionConfig = "yaml"
	PathConfig          = "./configs"
)

func Init(pathConfigs ...string) {
	pathConfig := PathConfig
	if len(pathConfigs) > 0 {
		pathConfig = pathConfigs[0]
	}

	viper.SetConfigName(FileNameConfig)
	viper.SetConfigType(FileExtensionConfig)
	viper.AddConfigPath(pathConfig)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read file config: %v", err))
	}

	for _, key := range viper.AllKeys() {
		value := viper.Get(key)
		viper.SetDefault(key, value)
	}
}
