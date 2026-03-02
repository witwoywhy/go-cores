package kafka

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/witwoywhy/go-cores/cryptos"
)

type FranzaGOConsumerGroup struct {
	Topic  string
	Group  string
	Client *kgo.Client

	ctx            context.Context
	MessageChannel chan *Message
}

func NewFranzaGOConsumerGroup(key string) *FranzaGOConsumerGroup {
	var config Config
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed when new consumer group unmarshal key %s: %v", key, err))
	}

	tls, err := cryptos.NewTLSConfig(config.Cert.CertFile, config.Cert.KeyFile, config.Cert.CaFile)
	if err != nil {
		panic(fmt.Errorf("failed when new consumer group tls config %s: %v", key, err))
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Broker),
		kgo.DialTLSConfig(tls),
		kgo.ConsumeTopics(config.Topic),
		kgo.ConsumerGroup(config.ConsumerGroup),
		kgo.AutoCommitMarks(),
	)

	consumer := &FranzaGOConsumerGroup{
		Topic:          config.Topic,
		Group:          config.ConsumerGroup,
		Client:         client,
		ctx:            context.Background(),
		MessageChannel: make(chan *Message),
	}

	go consumer.Consume()

	return consumer
}

func (c *FranzaGOConsumerGroup) Consume() {
	for {
		fetches := c.Client.PollFetches(c.ctx)
		if c.ctx.Err() != nil {
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "Fetch error — topic: %s, partition: %d, err: %v\n", e.Topic, e.Partition, e.Err)
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

func (c *FranzaGOConsumerGroup) Shutdown() error {
	c.Client.Close()
	return nil
}
