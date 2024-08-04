package gorms

import (
	"fmt"
	"strings"

	mysqld "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/witwoywhy/go-cores/dbs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(key string) *gorm.DB {
	var config dbs.DbConfig
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed to loaded config db %s: %v", key, err))
	}

	var gormDb *gorm.DB
	var err error
	switch strings.ToLower(config.Driver) {
	case dbs.Pg:
		gormDb, err = gorm.Open(postgres.New(postgres.Config{DSN: config.Dsn}))
	case dbs.Mysql:
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
		panic("failed to get db driver")
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
