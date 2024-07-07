package gorms

import (
	"fmt"
	"strings"
	"time"

	mysqld "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
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

func Init(key string) *gorm.DB {
	var config Config
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed to loaded config db %s: %v", key, err))
	}

	var gormDb *gorm.DB
	var err error
	switch strings.ToLower(config.Driver) {
	case "pg":
		gormDb, err = gorm.Open(postgres.New(postgres.Config{DSN: config.Dsn}))
	case "mysql":
		mysqlConfig := mysql.Config{}
		if config.Dsn != "" {
			mysqlConfig.DSN = config.Dsn
		} else {
			mysqlConfig.DSNConfig = &mysqld.Config{
				User:         config.Username,
				Passwd:       config.Password,
				Net:          "",
				Addr:         config.Host,
				DBName:       config.Database,
				Params:       map[string]string{},
				Timeout:      config.Timeout,
				ReadTimeout:  config.ReadTimeout,
				WriteTimeout: config.WriteTimeout,
			}
		}
		gormDb, err = gorm.Open(mysql.New(mysqlConfig))
	default:
		panic(fmt.Errorf("failed to get db driver: %v", err))
	}
	if err != nil {
		panic(fmt.Errorf("failed to open db %s: %v", key, err))
	}

	db, err := gormDb.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get db %s: %v", key, err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to ping db %s: %v", key, err))
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxConns)
	db.SetConnMaxLifetime(config.MaxLifeTime)

	return gormDb
}
