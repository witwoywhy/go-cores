package reqs

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/req"
)

type Client interface {
	Request(l logger.Logger) Request
	Config() *Config
}

type client struct {
	client *req.Client
	config *Config
}

func NewClient(key string) Client {
	var config Config
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed to new req client [%s]: %v", key, err))
	}

	c := client{
		client: req.NewClient(),
		config: &config,
	}

	if config.EnableInsecureSkipVerify {
		c.client.EnableInsecureSkipVerify()
	}

	c.client.BaseURL = config.BaseUrl
	if config.Timeout != 0 {
		c.client.SetTimeout(config.Timeout)
	}

	return c
}

func NewClientWithConfig(config *Config) Client {
	c := client{
		client: req.NewClient(),
		config: config,
	}

	if config.EnableInsecureSkipVerify {
		c.client.EnableInsecureSkipVerify()
	}

	c.client.BaseURL = config.BaseUrl
	if config.Timeout != 0 {
		c.client.SetTimeout(config.Timeout)
	}

	return &c
}

func (c client) Request(l logger.Logger) Request {
	return &request{
		request: c.client.NewRequest(),
		config:  c.config,
		l:       l,
	}
}

func (c client) Config() *Config {
	return c.config
}
