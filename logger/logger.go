package logger

type Logger interface {
	Info(obj any)
	Debug(obj any)
	Warn(obj any)
	Error(err any)

	Infof(format string, obj ...any)
	Debugf(format string, obj ...any)
	Warnf(format string, obj ...any)
	Errorf(format string, obj ...any)

	JSON(m map[string]any)
}
