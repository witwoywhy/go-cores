package circuitbreaker

import (
	"strings"

	"github.com/afex/hystrix-go/hystrix"
)

func Do(name string, fn func() error, fallBack func(error) error) error {
	return hystrix.Do(strings.ToLower(name), fn, fallBack)
}

func Go(name string, fn func() error, fallBack func(error) error) chan error {
	return hystrix.Go(name, fn, fallBack)
}
