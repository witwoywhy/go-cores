package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
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

	onceShutdown sync.Once
}

func NewProducer(options ...Option) pubsub.Producer {
	var config Config
	for _, option := range options {
		option.apply(&config)
	}

	if err := viper.UnmarshalKey(config.key, &config); err != nil {
		panic(fmt.Errorf("failed when kafka new producer viper.UnmarshalKey [%s]: %v", config.key, err))
	}

	if config.Broker == "" {
		fmt.Println(`new producer config.Broker is ""`)
		return nil
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
		kgo.DefaultProduceTopic(config.Topic),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)

	if err := client.Ping(context.Background()); err != nil {
		panic(fmt.Errorf("failed when kafka producer ping [%s]: %v", config.key, err))
	}

	return &Producer{
		Topic:        config.Topic,
		Client:       client,
		onceShutdown: sync.Once{},
	}
}

func (p *Producer) Produce(key string, v any, l logger.Logger) error {
	var b []byte
	switch t := v.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		var err error
		b, err = json.Marshal(v)
		if err != nil {
			return err
		}
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
	p.onceShutdown.Do(func() {
		p.Client.Close()
	})
	return nil
}
