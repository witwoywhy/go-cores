package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/cryptos"
)

type Producer struct {
	Topic  string
	Client sarama.SyncProducer
}

func NewProducer(key string) *Producer {
	var config Config
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed when new producer unmarshal key %s: %v", key, err))
	}

	tls, err := cryptos.NewTLSConfig(config.Cert.CertFile, config.Cert.KeyFile, config.Cert.CaFile)
	if err != nil {
		panic(fmt.Errorf("failed when new producer tls config %s: %v", key, err))
	}

	cfg := sarama.NewConfig()
	cfg.Net.TLS.Enable = true
	cfg.Net.TLS.Config = tls
	cfg.Version = sarama.V3_0_0_0
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{config.Broker}, cfg)
	if err != nil {
		panic(fmt.Errorf("failed when create producer %s: %v", key, err))
	}

	return &Producer{
		Topic:  config.Topic,
		Client: producer,
	}
}

func (p *Producer) Publish(key string, v any) error {
	var b sarama.ByteEncoder
	switch t := v.(type) {
	case string:
		b = sarama.ByteEncoder(string(t))
	default:
		bb, err := json.Marshal(v)
		if err != nil {
			return err
		}

		b = bb
	}

	_, _, err := p.Client.SendMessage(&sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(key),
		Value: b,
	})

	return err
}

func (p *Producer) Shutdown() error {
	return p.Client.Close()
}
