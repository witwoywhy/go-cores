package pubsub

import "github.com/witwoywhy/go-cores/logger"

type Publisher interface {
	Publish(key string, v any, l logger.Logger) error
	Shutdown(l logger.Logger) error
}
