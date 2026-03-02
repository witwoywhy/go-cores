package kafka

import "github.com/witwoywhy/go-cores/logger"

type Option interface{ apply(*Config) }

type clientOption struct{ fn func(*Config) }

func (opt clientOption) apply(config *Config) { opt.fn(config) }

type Config struct {
	Broker        string     `mapstructure:"broker"`
	Topic         string     `mapstructure:"topic"`
	ConsumerGroup string     `mapstructure:"consumerGroup"`
	Cert          CertConfig `mapstructure:"cert"`

	key string
	l   logger.Logger
}

// type ProduceConfig struct {
// 	MaxMessageBytes int           `mapstructure:"maxMessageBytes"`
// 	MaxRetry        int           `mapstructure:"maxRetry"`
// 	RetryBackoff    time.Duration `mapstructure:"retryBackOff"`
// }

// type ConsumeConfig struct {
// 	MaxFetch    int           `mapstructure:"maxFetch"`
// 	MaxWaitTime time.Duration `mapstructure:"maxWaitTime"`
// }

type CertConfig struct {
	CertFile string `mapstructure:"certFile"`
	KeyFile  string `mapstructure:"keyFile"`
	CaFile   string `mapstructure:"caFile"`
}

func Key(key string) Option {
	return clientOption{func(c *Config) { c.key = key }}
}

func Log(l logger.Logger) Option {
	return clientOption{func(c *Config) { c.l = l }}
}
