package dbs

import (
	"testing"

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
