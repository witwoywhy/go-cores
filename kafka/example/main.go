package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/witwoywhy/go-cores/kafka"
	"github.com/witwoywhy/go-cores/logs"
	"github.com/witwoywhy/go-cores/vipers"
)

func init() {
	vipers.Init()
}

func main() {
	c := kafka.NewConsumerGroup(
		kafka.AddConfigKey("pubsub.sub"),
		kafka.AddLogger(logs.L),
	)

	var wg sync.WaitGroup
	wg.Go(func() {
		c.Consume(logs.L, func(ctx context.Context, topic, group, key string, value []byte) bool {
			fmt.Println("")
			fmt.Println(string(value))
			fmt.Println("")
			return true
		}, func() {})
	})

	time.Sleep(5 * time.Second)
	p := kafka.NewProducer(
		kafka.AddConfigKey("pubsub.pub"),
		kafka.AddLogger(logs.L),
	)
	fmt.Println(p.Produce("AAA", "HELLO FROM PUB", logs.L))

	wg.Wait()
}
