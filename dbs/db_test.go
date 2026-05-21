package dbs

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDbConfigToDSN(t *testing.T) {
	t.Run("pg uses explicit dsn", func(t *testing.T) {
		config := &DbConfig{
			Driver: Pg,
			DSN:    "postgres://root:root@localhost:5432/witwoywhy",
		}

		assert.Equal(t, "postgres://root:root@localhost:5432/witwoywhy", config.ToDSN())
	})

	t.Run("postgres driver uses explicit dsn", func(t *testing.T) {
		config := &DbConfig{
			Driver: Postgres,
			DSN:    "postgres://root:root@localhost:5432/witwoywhy",
		}

		assert.Equal(t, "postgres://root:root@localhost:5432/witwoywhy", config.ToDSN())
	})

	t.Run("pg builds dsn from fields", func(t *testing.T) {
		config := &DbConfig{
			Driver:   Pg,
			Username: "root",
			Password: "root",
			Host:     "localhost",
			Port:     "5432",
			Database: "witwoywhy",
		}

		assert.Equal(t, "postgres://root:root@localhost:5432/witwoywhy", config.ToDSN())
	})

	t.Run("pg builds dsn from fields with connect timeout", func(t *testing.T) {
		config := &DbConfig{
			Driver:   Pg,
			Username: "root",
			Password: "root",
			Host:     "localhost",
			Port:     "5432",
			Database: "witwoywhy",
			Timeout:  5 * time.Second,
		}

		dsn, err := url.Parse(config.ToDSN())

		assert.NoError(t, err)
		assert.Equal(t, "postgres", dsn.Scheme)
		assert.Equal(t, "root:root@localhost:5432", dsn.User.String()+"@"+dsn.Host)
		assert.Equal(t, "/witwoywhy", dsn.Path)
		assert.Equal(t, "5", dsn.Query().Get("connect_timeout"))
	})

	t.Run("pg uses explicit url dsn with connect timeout from config", func(t *testing.T) {
		config := &DbConfig{
			Driver:  Pg,
			DSN:     "postgres://root:root@localhost:5432/witwoywhy",
			Timeout: 5 * time.Second,
		}

		dsn, err := url.Parse(config.ToDSN())

		assert.NoError(t, err)
		assert.Equal(t, "5", dsn.Query().Get("connect_timeout"))
	})

	t.Run("pg keeps explicit connect timeout", func(t *testing.T) {
		config := &DbConfig{
			Driver:  Pg,
			DSN:     "postgres://root:root@localhost:5432/witwoywhy?connect_timeout=3",
			Timeout: 5 * time.Second,
		}

		assert.Equal(t, "postgres://root:root@localhost:5432/witwoywhy?connect_timeout=3", config.ToDSN())
	})

	t.Run("pg preserves explicit url query params with connect timeout", func(t *testing.T) {
		config := &DbConfig{
			Driver:  Pg,
			DSN:     "postgres://root:root@localhost:5432/witwoywhy?sslmode=disable",
			Timeout: 1500 * time.Millisecond,
		}

		dsn, err := url.Parse(config.ToDSN())

		assert.NoError(t, err)
		assert.Equal(t, "disable", dsn.Query().Get("sslmode"))
		assert.Equal(t, "2", dsn.Query().Get("connect_timeout"))
	})

	t.Run("pg keeps explicit non url dsn unchanged", func(t *testing.T) {
		config := &DbConfig{
			Driver:  Pg,
			DSN:     "host=localhost user=root password=root dbname=witwoywhy",
			Timeout: 5 * time.Second,
		}

		assert.Equal(t, "host=localhost user=root password=root dbname=witwoywhy", config.ToDSN())
	})

	t.Run("mysql uses explicit dsn", func(t *testing.T) {
		config := &DbConfig{
			Driver: Mysql,
			DSN:    "root:root@tcp(localhost:3306)/witwoywhy?parseTime=true",
		}

		assert.Equal(t, "root:root@tcp(localhost:3306)/witwoywhy?parseTime=true", config.ToDSN())
	})

	t.Run("mysql builds dsn from fields", func(t *testing.T) {
		config := &DbConfig{
			Driver:   Mysql,
			Username: "root",
			Password: "root",
			Host:     "localhost",
			Port:     "3306",
			Database: "witwoywhy",
		}

		assert.Equal(t, "root:root@tcp(localhost:3306)/witwoywhy", config.ToDSN())
	})
}
