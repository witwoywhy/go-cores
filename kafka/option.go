package kafka

import "github.com/witwoywhy/go-cores/logger"

type Option interface{ apply(*Config) }

type clientOption struct{ fn func(*Config) }

func (opt clientOption) apply(config *Config) { opt.fn(config) }

func AddConfigKey(key string) Option {
	return clientOption{func(c *Config) { c.key = key }}
}

func AddLogger(l logger.Logger) Option {
	return clientOption{func(c *Config) { c.l = l }}
}
