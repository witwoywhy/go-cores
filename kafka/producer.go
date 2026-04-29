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
	var (
		opt    OptionConfig
		config ProducerConfig
	)
	for _, option := range options {
		option.apply(&opt)
	}

	if err := viper.UnmarshalKey(opt.key, &config); err != nil {
		panic(fmt.Errorf("failed when kafka new producer viper.UnmarshalKey [%s]: %v", opt.key, err))
	}

	config = checkingProducerConfig(config)

	var tls *tls.Config
	var err error
	switch config.Cert.Type {
	case "value":
		tls, err = cryptos.NewTLSConfig(config.Cert.Cert, config.Cert.Key, config.Cert.CA)
	default: // flie
		tls, err = cryptos.NewTLSConfigFromFile(config.Cert.Cert, config.Cert.Key, config.Cert.CA)
	}
	if err != nil {
		panic(fmt.Errorf("failed when kafka new consumer group tls config [%s]: %v", opt.key, err))
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Broker),
		kgo.DialTLSConfig(tls),
		kgo.DefaultProduceTopic(config.Topic),

		kgo.ProducerBatchMaxBytes(int32(config.BatchSize)),
		kgo.MaxBufferedRecords(config.MaxBufferRecords),
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),

		kgo.RequiredAcks(kgo.AllISRAcks()),

		kgo.ProduceRequestTimeout(config.RequestTimeout),
		kgo.RecordDeliveryTimeout(config.DeliveryTimeout),
		kgo.RetryBackoffFn(jitteredBackoff),
	)
	if err := client.Ping(context.Background()); err != nil {
		panic(fmt.Errorf("failed when kafka producer ping [%s]: %v", opt.key, err))
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

var defaultProducerConfig = ProducerConfig{
	BatchSize:        1048576,
	MaxBufferRecords: 10000,
	Linger:           10 * time.Millisecond,
	RequestTimeout:   60 * time.Second,
	DeliveryTimeout:  180 * time.Second,
}

func checkingProducerConfig(config ProducerConfig) ProducerConfig {
	var cfg = &defaultProducerConfig
	cfg.Broker = config.Broker
	cfg.Topic = config.Topic
	cfg.Cert = config.Cert

	if config.BatchSize > 0 {
		cfg.BatchSize = config.BatchSize
	}

	if config.MaxBufferRecords > 0 {
		cfg.MaxBufferRecords = config.MaxBufferRecords
	}

	if config.Linger > 0 {
		cfg.Linger = config.Linger
	}

	if config.RequestTimeout > 0 {
		cfg.RequestTimeout = config.RequestTimeout
	}

	if config.DeliveryTimeout > 0 {
		cfg.DeliveryTimeout = config.DeliveryTimeout
	}

	return *cfg
}
