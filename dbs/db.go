package dbs

import (
	"fmt"
	"time"
)

type DbConfig struct {
	Driver       string        `mapstructure:"driver"`
	DSN          string        `mapstructure:"dsn"`
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	Database     string        `mapstructure:"database"`
	Timeout      time.Duration `mapstructure:"timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxConns     int           `mapstructure:"max_conns"`
	MaxLifeTime  time.Duration `mapstructure:"max_life_time"`
}

func (d *DbConfig) ToDSN() string {
	switch d.Driver {
	case Pg:
		if d.DSN == "" {
			return d.DSN
		}

		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
	case Mysql:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
	default:
		panic("failed to get db DSN")
	}
}
