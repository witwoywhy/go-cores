package gorms

import (
	"context"
	"time"

	ll "github.com/witwoywhy/go-cores/logger"
	"gorm.io/gorm/logger"
)

type gormLog struct {
	l ll.Logger
}

func NewGormLog(l ll.Logger) logger.Interface {
	return &gormLog{
		l: l,
	}
}

func (g *gormLog) Error(context.Context, string, ...interface{}) {
	panic("unimplemented gorm log error")
}

func (g *gormLog) Info(context.Context, string, ...interface{}) {
	panic("unimplemented gorm log info")
}

func (g *gormLog) LogMode(logger.LogLevel) logger.Interface {
	panic("unimplemented gorm log mode")
}

func (g *gormLog) Warn(context.Context, string, ...interface{}) {
	panic("unimplemented gorm log warn")
}

func (g *gormLog) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	g.l.Info(StartOutbound)
	sql, _ := fc()
	since := time.Since(begin)
	g.l.Infof(EndOutbound, since, getTableNameFromQuery(sql))
}
