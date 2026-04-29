package kafka

import (
	"time"

	"github.com/witwoywhy/go-cores/logger"
)

type OptionConfig struct {
	key string
	l   logger.Logger
}

type ProducerConfig struct {
	Broker string     `mapstructure:"broker"`
	Topic  string     `mapstructure:"topic"`
	Cert   CertConfig `mapstructure:"cert"`

	// Message Size
	// batch.size | default=1048576 (1MB)
	BatchSize int `mapstructure:"batch_size"`
	// default=10000
	MaxBufferRecords int `mapstructure:"max_buffer_records"`

	// Timeouts
	// linger.ms | default=10ms
	Linger time.Duration `mapstructure:"linger"`
	// request.timeout.ms | default=60s
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	// delivery.timeout.ms | default=180s
	DeliveryTimeout time.Duration `mapstructure:"delivery_timeout"`
}

type ConsumerConfig struct {
	Broker        string     `mapstructure:"broker"`
	Topic         string     `mapstructure:"topic"`
	ConsumerGroup string     `mapstructure:"consumer_group"`
	Cert          CertConfig `mapstructure:"cert"`

	// Message Size
	// fetch.max.bytes | defualt=5242880(5MB)
	FetchMaxBytes int `mapstructure:"fetch_max_bytes"`
	// fetch.min.bytes | defualt=1024
	FetchMinBytes int `mapstructure:"fetch_min_bytes"`
	// max.partition.fetch.bytes | default=1048576 (1MB)
	FetchMaxParitionBytes int `mapstructure:"fetch_max_partition_bytes"`

	// fetch.max.wait.ms | default=1s
	FetchMaxWait time.Duration `mapstructure:"fetch_max_wait"`
	// session.timeout.ms |  default=60s
	SessionTimeout time.Duration `mapstructure:"session_timeout"`
	// heartbeat.interval.ms | default=10s
	HeartbeatInternal time.Duration `mapstructure:"heartbeat_interval"`
	// max.poll.interval.ms | default=600s
	MaxPollInterval time.Duration `mapstructure:"max_poll_interval"`
	// request.timeout.ms | default=60s
	RequestTimeout time.Duration `mapstructure:"request_timeout"`

	// consume.reset.offset | default=end (start, end)
	ConsumeResetOffset string `mapstructure:"consume_reset_offset"`
}

type CertConfig struct {
	// file or value
	Type string `mapstructure:"type"`

	CA   string `mapstructure:"ca"`
	Key  string `mapstructure:"key"`
	Cert string `mapstructure:"cert"`
}
