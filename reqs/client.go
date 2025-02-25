package reqs

import (
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

type Client interface {
	Request() Request
}

type client struct {
	client *req.Client
	config *Config
}

func NewClient(key string) Client {
	var config Config
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed to new req client: %v", err))
	}

	c := client{
		client: req.NewClient(),
		config: &config,
	}

	c.client.BaseURL = config.BaseUrl
	if config.Timeout != 0 {
		c.client.SetTimeout(config.Timeout)
	}

	return &c
}

func NewClientWithConfig(config Config) Client {
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

	return &c
}

func (c *client) Request() Request {
	request := request{
		request: c.client.NewRequest(),
		config:  c.config,
	}
	return &request
}
