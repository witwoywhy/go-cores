package kafka

import "github.com/witwoywhy/go-cores/logger"

type Option interface{ apply(*OptionConfig) }

type clientOption struct{ fn func(*OptionConfig) }

func (opt clientOption) apply(config *OptionConfig) { opt.fn(config) }

func AddConfigKey(key string) Option {
	return clientOption{func(c *OptionConfig) { c.key = key }}
}

func AddLogger(l logger.Logger) Option {
	return clientOption{func(c *OptionConfig) { c.l = l }}
}
