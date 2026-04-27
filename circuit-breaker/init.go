package circuitbreaker

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/spf13/viper"
)

func Init() {
	configs := viper.GetStringMap("circuit_breakers")

	if len(configs) == 0 {
		return
	}

	for k, v := range configs {
		hystrix.ConfigureCommand(k, hystrix.CommandConfig{
			Timeout:                v.(map[string]any)["timeout"].(int),
			MaxConcurrentRequests:  v.(map[string]any)["max_concurrent_requests"].(int),
			RequestVolumeThreshold: v.(map[string]any)["request_volume_threshold"].(int),
			SleepWindow:            v.(map[string]any)["sleep_window"].(int),
			ErrorPercentThreshold:  v.(map[string]any)["error_percent_threshold"].(int),
		})
	}
}
