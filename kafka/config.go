package kafka

import "github.com/witwoywhy/go-cores/logger"

type Config struct {
	Broker        string     `mapstructure:"broker"`
	Topic         string     `mapstructure:"topic"`
	ConsumerGroup string     `mapstructure:"consumer_group"`
	Cert          CertConfig `mapstructure:"cert"`

	key string
	l   logger.Logger
}

type CertConfig struct {
	// file or value
	Type string `mapstructure:"type"`

	CA   string `mapstructure:"ca"`
	Key  string `mapstructure:"key"`
	Cert string `mapstructure:"cert"`
}
