package pubsub

import "context"

// return bool, true=ack
type HandlerFunc = func(ctx context.Context, topic, group, key string, value []byte) bool
