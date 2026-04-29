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
	var (
		opt    OptionConfig
		config ConsumerConfig
	)
	for _, option := range options {
		option.apply(&opt)
	}

	if err := viper.UnmarshalKey(opt.key, &config); err != nil {
		panic(fmt.Errorf("failed when kafka new consumer group viper.UnmarshalKey [%s]: %v", opt.key, err))
	}

	config = checkingConsumerConfig(config)

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
		kgo.ConsumeTopics(config.Topic),
		kgo.ConsumerGroup(config.ConsumerGroup),

		kgo.FetchMaxBytes(int32(config.FetchMaxBytes)),
		kgo.FetchMinBytes(int32(config.FetchMinBytes)),
		kgo.FetchMaxPartitionBytes(int32(config.FetchMaxParitionBytes)),

		kgo.FetchMaxWait(config.FetchMaxWait),
		kgo.SessionTimeout(config.SessionTimeout),
		kgo.HeartbeatInterval(config.HeartbeatInternal),
		kgo.RebalanceTimeout(config.MaxPollInterval),
		kgo.RequestTimeoutOverhead(config.RequestTimeout),

		kgo.ConsumeResetOffset(config.consumeResetOffset()),
		kgo.FetchIsolationLevel(kgo.ReadCommitted()),
	)
	if err := client.Ping(context.Background()); err != nil {
		panic(fmt.Errorf("failed when kafka consumer group ping [%s]: %v", opt.key, err))
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

var defaultConsumerConfig = ConsumerConfig{
	FetchMaxBytes:         0,
	FetchMinBytes:         0,
	FetchMaxParitionBytes: 0,
	FetchMaxWait:          0,
	SessionTimeout:        0,
	HeartbeatInternal:     0,
	MaxPollInterval:       0,
	RequestTimeout:        0,
	ConsumeResetOffset:    "end",
}

func checkingConsumerConfig(config ConsumerConfig) ConsumerConfig {
	cfg := defaultConsumerConfig
	cfg.Broker = config.Broker
	cfg.Topic = config.Topic
	cfg.ConsumerGroup = config.ConsumerGroup
	cfg.Cert = config.Cert

	if config.FetchMaxBytes > 0 {
		cfg.FetchMaxBytes = config.FetchMaxBytes
	}

	if config.FetchMinBytes > 0 {
		cfg.FetchMinBytes = config.FetchMinBytes
	}

	if config.FetchMaxParitionBytes > 0 {
		cfg.FetchMaxParitionBytes = config.FetchMaxParitionBytes
	}

	if config.FetchMaxWait > 0 {
		cfg.FetchMaxWait = config.FetchMaxWait
	}

	if config.SessionTimeout > 0 {
		cfg.SessionTimeout = config.SessionTimeout
	}

	if config.HeartbeatInternal > 0 {
		cfg.HeartbeatInternal = config.HeartbeatInternal
	}

	if config.MaxPollInterval > 0 {
		cfg.MaxPollInterval = config.MaxPollInterval
	}

	if config.RequestTimeout > 0 {
		cfg.RequestTimeout = config.RequestTimeout
	}

	if config.ConsumeResetOffset != "" {
		cfg.ConsumeResetOffset = config.ConsumeResetOffset
	}

	return cfg
}

func (c ConsumerConfig) consumeResetOffset() kgo.Offset {
	switch c.ConsumeResetOffset {
	case "start":
		return kgo.NewOffset().AtStart()
	default:
		return kgo.NewOffset().AtEnd()
	}
}
