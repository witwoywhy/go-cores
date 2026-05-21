package gorms

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type pingDB struct {
	pingCalled         bool
	pingContextCalled  bool
	contextHasDeadline bool
}

func (p *pingDB) Ping() error {
	p.pingCalled = true
	return nil
}

func (p *pingDB) PingContext(ctx context.Context) error {
	p.pingContextCalled = true
	_, p.contextHasDeadline = ctx.Deadline()
	return nil
}

func TestPing(t *testing.T) {
	t.Run("uses ping when timeout is missing", func(t *testing.T) {
		db := &pingDB{}

		err := ping(db, 0)

		assert.NoError(t, err)
		assert.True(t, db.pingCalled)
		assert.False(t, db.pingContextCalled)
	})

	t.Run("uses ping context when timeout is configured", func(t *testing.T) {
		db := &pingDB{}

		err := ping(db, 5*time.Second)

		assert.NoError(t, err)
		assert.False(t, db.pingCalled)
		assert.True(t, db.pingContextCalled)
		assert.True(t, db.contextHasDeadline)
	})
}
