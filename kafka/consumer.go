package kafka

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/witwoywhy/go-cores/cryptos"
	"github.com/witwoywhy/go-cores/logger"
)

type ConsumerGroup struct {
	Topic  string
	Group  string
	Client *kgo.Client

	ctx            context.Context
	l              logger.Logger
	MessageChannel chan *Message
}

func NewConsumerGroup(options ...Option) *ConsumerGroup {
	var config Config
	for _, option := range options {
		option.apply(&config)
	}

	if err := viper.UnmarshalKey(config.key, &config); err != nil {
		panic(fmt.Errorf("failed when new consumer group unmarshal key %s: %v", config.key, err))
	}

	tls, err := cryptos.NewTLSConfig(config.Cert.CertFile, config.Cert.KeyFile, config.Cert.CaFile)
	if err != nil {
		panic(fmt.Errorf("failed when new consumer group tls config %s: %v", config.key, err))
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Broker),
		kgo.DialTLSConfig(tls),
		kgo.ConsumeTopics(config.Topic),
		kgo.ConsumerGroup(config.ConsumerGroup),
		kgo.AutoCommitMarks(),
	)

	consumer := &ConsumerGroup{
		Topic:          config.Topic,
		Group:          config.ConsumerGroup,
		Client:         client,
		ctx:            context.Background(),
		MessageChannel: make(chan *Message),
	}

	go consumer.Consume(config.l)

	return consumer
}

func (c *ConsumerGroup) Consume(l logger.Logger) {
	l.Infof("START CONSUME %v", c.Group)
	for {
		fetches := c.Client.PollFetches(c.ctx)
		if c.ctx.Err() != nil {
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, e := range errs {
				c.l.Errorf("fetch error — topic: %s, partition: %d, err: %v\n", e.Topic, e.Partition, e.Err)
			}
			continue
		}

		fetches.EachRecord(func(r *kgo.Record) {
			c.MessageChannel <- &Message{
				Key:       string(r.Key),
				Value:     string(r.Value),
				Topic:     r.Topic,
				Partition: r.Partition,
				Offset:    r.Offset,
			}
		})
	}
}

func (c *ConsumerGroup) Shutdown() error {
	c.Client.Close()
	return nil
}
