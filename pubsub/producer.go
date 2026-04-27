package pubsub

import "github.com/witwoywhy/go-cores/logger"

type Producer interface {
	Produce(key string, v any, l logger.Logger) error
	Shutdown(l logger.Logger) error
}
