package dbs

import "time"

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
