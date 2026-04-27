package gorms

import (
	"context"
	"fmt"
	"time"

	ll "github.com/witwoywhy/go-cores/logger"
)

type GormLogConfigOption interface{ apply(*gormLog) }

type clientGormLogConfigOption struct{ fn func(*gormLog) }

func (opt clientGormLogConfigOption) apply(config *gormLog) { opt.fn(config) }

func AddLogger(l ll.Logger) GormLogConfigOption {
	return clientGormLogConfigOption{
		fn: func(config *gormLog) {
			config.l = l
		},
	}
}

func Debug() GormLogConfigOption {
	return clientGormLogConfigOption{
		fn: func(config *gormLog) {
			config.traceFunc = traceWithDebug
		},
	}
}

func traceWithoutDebug(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error, l ll.Logger) {
	l.Info(StartOutbound)
	sql, _ := fc()
	since := time.Since(begin)
	l.Infof(EndOutbound, since, getTableNameFromQuery(sql))
}

func traceWithDebug(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error, l ll.Logger) {
	sql, _ := fc()
	since := time.Since(begin)
	fmt.Println("")
	fmt.Printf("[%v] => %s", since, sql)
	fmt.Println("")
}
