package circuitbreaker

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/spf13/viper"
)

func Init() {
	configs := viper.GetStringMap("circuitBreakers")

	if len(configs) == 0 {
		return
	}

	for k, v := range configs {
		hystrix.ConfigureCommand(k, hystrix.CommandConfig{
			Timeout:                v.(map[string]any)["timeout"].(int),
			MaxConcurrentRequests:  v.(map[string]any)["maxconcurrentrequests"].(int),
			RequestVolumeThreshold: v.(map[string]any)["requestvolumethreshold"].(int),
			SleepWindow:            v.(map[string]any)["sleepwindow"].(int),
			ErrorPercentThreshold:  v.(map[string]any)["errorpercentthreshold"].(int),
		})
	}
}
