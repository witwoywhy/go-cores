package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/witwoywhy/go-cores/cryptos"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/pubsub"
)

type Producer struct {
	Topic  string
	Client *kgo.Client
}

func NewProducer(options ...Option) pubsub.Publisher {
	var config Config
	for _, option := range options {
		option.apply(&config)
	}

	if err := viper.UnmarshalKey(config.key, &config); err != nil {
		panic(fmt.Errorf("failed when NewFranzaGOProducer viper.UnmarshalKey %s: %v", config.key, err))
	}

	tls, err := cryptos.NewTLSConfig(config.Cert.CertFile, config.Cert.KeyFile, config.Cert.CaFile)
	if err != nil {
		panic(fmt.Errorf("failed when new producer tls config %s: %v", config.key, err))
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Broker),
		kgo.DialTLSConfig(tls),
		kgo.DefaultProduceTopic(config.Topic),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)

	return &Producer{
		Topic:  config.Topic,
		Client: client,
	}
}

func (p *Producer) Publish(key string, v any, l logger.Logger) error {
	var b []byte
	switch t := v.(type) {
	case string:
		b = []byte(t)
	default:
		marshalByte, err := json.Marshal(v)
		if err != nil {
			return err
		}

		b = marshalByte
	}

	p.Client.Produce(context.Background(), &kgo.Record{
		Key:       []byte(key),
		Value:     b,
		Timestamp: time.Now(),
		Topic:     p.Topic,
	}, func(r *kgo.Record, err error) {
		if err != nil {
			l.Errorf("failed when produce %s: %v", key, err)
			return
		}
	})

	return nil
}

func (p *Producer) Shutdown(l logger.Logger) error {
	p.Client.Close()
	return nil
}
