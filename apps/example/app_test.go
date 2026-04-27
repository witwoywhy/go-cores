package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/vipers"
)

func TestInit(t *testing.T) {
	vipers.Init()
	apps.Init()

	var config = apps.ConfigInfo{
		Name:     "example",
		Env:      "dev",
		TimeZone: "UTC",
	}

	assert.Equal(t, config, apps.Config)
}
