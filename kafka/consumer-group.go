package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/cryptos"
)

type ConsumerGroup struct {
	Topic          string
	Group          string
	Client         sarama.ConsumerGroup
	MessageChannel chan *Message
}

func NewConsumerGroup(key string) *ConsumerGroup {
	var config Config
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed when new consumer group unmarshal key %s: %v", key, err))
	}

	tls, err := cryptos.NewTLSConfig(config.Cert.CertFile, config.Cert.KeyFile, config.Cert.CaFile)
	if err != nil {
		panic(fmt.Errorf("failed when new consumer group tls config %s: %v", key, err))
	}

	cfg := sarama.NewConfig()
	cfg.Net.TLS.Enable = true
	cfg.Net.TLS.Config = tls
	cfg.Version = sarama.V3_0_0_0
	cfg.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumerGroup([]string{config.Broker}, config.ConsumerGroup, cfg)
	if err != nil {
		panic(fmt.Errorf("failed when create consumer group %s: %v", key, err))
	}

	var (
		messageChannel = make(chan *Message)
		setupChannel   = make(chan struct{})
	)

	go consumer.Consume(
		context.Background(),
		[]string{config.Topic},
		NewConsumerGroupHandler(
			messageChannel,
			setupChannel,
		),
	)

	<-setupChannel
	close(setupChannel)

	return &ConsumerGroup{
		Topic:          config.Topic,
		Group:          config.ConsumerGroup,
		Client:         consumer,
		MessageChannel: messageChannel,
	}
}

func (c *ConsumerGroup) Shutdown() error {
	return c.Client.Close()
}

