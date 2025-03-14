package vipers

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName(fileNameConfig)
	viper.SetConfigType(fileExtensionConfig)
	viper.AddConfigPath(pathConfig)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read file config: %v", err))
	}

	for _, key := range viper.AllKeys() {
		value := viper.Get(key)
		viper.Set(key, value)
	}
}
