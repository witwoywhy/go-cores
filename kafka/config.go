package kafka

type Config struct {
	Broker        string     `mapstructure:"broker"`
	Topic         string     `mapstructure:"topic"`
	ConsumerGroup string     `mapstructure:"consumerGroup"`
	Cert          CertConfig `mapstructure:"cert"`
}

// type ProduceConfig struct {
// 	MaxMessageBytes int           `mapstructure:"maxMessageBytes"`
// 	MaxRetry        int           `mapstructure:"maxRetry"`
// 	RetryBackoff    time.Duration `mapstructure:"retryBackOff"`
// }

// type ConsumeConfig struct {
// 	MaxFetch    int           `mapstructure:"maxFetch"`
// 	MaxWaitTime time.Duration `mapstructure:"maxWaitTime"`
// }

type CertConfig struct {
	CertFile string `mapstructure:"certFile"`
	KeyFile  string `mapstructure:"keyFile"`
	CaFile   string `mapstructure:"caFile"`
}
