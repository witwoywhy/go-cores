package logs

import "github.com/witwoywhy/go-cores/pubsub"

type logConfigOption struct {
	producer pubsub.Producer
}

type LogConfigOption interface{ apply(*logConfigOption) }

type clientLogConfigOption struct{ fn func(*logConfigOption) }

func (opt clientLogConfigOption) apply(config *logConfigOption) { opt.fn(config) }

func AddProducer(producer pubsub.Producer) LogConfigOption {
	return clientLogConfigOption{
		fn: func(c *logConfigOption) {
			c.producer = producer
		},
	}
}
