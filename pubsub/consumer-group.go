package pubsub

import "github.com/witwoywhy/go-cores/logger"

type ConsumerGroup interface {
	Consume(l logger.Logger, fn HandlerFunc, closeFunc func())
}
