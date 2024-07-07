package logger

type Logger interface {
	Info(obj any)
	Debug(obj any)
	Warn(obj any)
	Error(err error)

	Infof(format string, obj ...any)
	Debugf(format string, obj ...any)
	Warnf(format string, obj ...any)
	Errorf(format string, obj ...any)
}

type CoreLogger interface {
	Logger

	JSON(m map[string]any)
	AddInformation(m map[string]any)
}
