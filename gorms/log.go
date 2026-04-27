package gorms

import (
	"context"
	"time"

	ll "github.com/witwoywhy/go-cores/logger"
	"gorm.io/gorm/logger"
)

type gormLog struct {
	l         ll.Logger
	traceFunc func(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error, l ll.Logger)
}

type GormLogOption struct {
	IsLogSql bool
}

func NewGormLog(options ...GormLogConfigOption) logger.Interface {
	var gl gormLog
	gl.traceFunc = traceWithoutDebug

	for _, option := range options {
		option.apply(&gl)
	}

	return &gl
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

func (g *gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	g.traceFunc(ctx, begin, fc, err, g.l)
}
