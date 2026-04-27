package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/witwoywhy/go-cores/logger"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/witwoywhy/go-cores/cryptos"
	"github.com/witwoywhy/go-cores/pubsub"
)

type consumerGroup struct {
	Topic  string
	Group  string
	Client *kgo.Client

	ctx context.Context

	onceShutdown sync.Once
}

func NewConsumerGroup(options ...Option) pubsub.ConsumerGroup {
	var config Config
	for _, option := range options {
		option.apply(&config)
	}

	if err := viper.UnmarshalKey(config.key, &config); err != nil {
		panic(fmt.Errorf("failed when kafka new consumer group viper.UnmarshalKey [%s]: %v", config.key, err))
	}

	var tls *tls.Config
	var err error
	switch config.Cert.Type {
	case "value":
		tls, err = cryptos.NewTLSConfig(config.Cert.Cert, config.Cert.Key, config.Cert.CA)
	default: // flie
		tls, err = cryptos.NewTLSConfigFromFile(config.Cert.Cert, config.Cert.Key, config.Cert.CA)
	}
	if err != nil {
		panic(fmt.Errorf("failed when kafka new consumer group tls config [%s]: %v", config.key, err))
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Broker),
		kgo.DialTLSConfig(tls),
		kgo.ConsumeTopics(config.Topic),
		kgo.ConsumerGroup(config.ConsumerGroup),
	)

	if err := client.Ping(context.Background()); err != nil {
		panic(fmt.Errorf("failed when kafka consumer group ping [%s]: %v", config.key, err))
	}

	return &consumerGroup{
		Topic:  config.Topic,
		Group:  config.ConsumerGroup,
		Client: client,
		ctx:    context.Background(),
	}
}

func (c *consumerGroup) Consume(l logger.Logger, fn pubsub.HandlerFunc, closeFunc func()) {
	go func() {
		l.Infof("START CONSUME: TOPIC [%v], GROUP [%v]", c.Topic, c.Group)
		for {
			fetches := c.Client.PollFetches(c.ctx)
			if c.ctx.Err() != nil {
				return
			}

			if errs := fetches.Errors(); len(errs) > 0 {
				for _, e := range errs {
					if !errors.Is(e.Err, kgo.ErrClientClosed) {
						l.Errorf("fetch error — topic: %s, partition: %d, err: %v\n", e.Topic, e.Partition, e.Err)
					}
				}
				continue
			}

			fetches.EachRecord(func(r *kgo.Record) {
				if fn(r.Context, c.Topic, c.Group, string(r.Key), r.Value) {
					c.Client.MarkCommitRecords(r)
				}
			})
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Info("start shutdown consumer ...")

	c.onceShutdown.Do(func() {
		c.Client.Close()

		if closeFunc != nil {
			closeFunc()
		}
	})

	l.Info("consumer shutdown")
}
