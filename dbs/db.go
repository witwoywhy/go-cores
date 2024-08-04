package dbs

import (
	"fmt"
	"time"
)

type DbConfig struct {
	Driver       string        `mapstructure:"driver"`
	Dsn          string        `mapstructure:"dsn"`
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	Database     string        `mapstructure:"database"`
	Timeout      time.Duration `mapstructure:"timeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	MaxIdleConns int           `mapstructure:"maxIdleConns"`
	MaxConns     int           `mapstructure:"maxConns"`
	MaxLifeTime  time.Duration `mapstructure:"maxLifeTime"`
}

func (d *DbConfig) ToDsn() string {
	switch d.Driver {
	case Pg:
		return d.Dsn
	case Mysql:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
	default:
		panic("failed to get db driver")
	}
}
