package logger

type Writer interface {
	Write(fields map[string]any, l Logger)
}
