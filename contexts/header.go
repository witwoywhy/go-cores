package contexts

import "github.com/witwoywhy/go-cores/enum/language"

type Header struct {
	Authorization string `header:"Authorization"`

	RequestID string `header:"X-Request-Id"`
	TraceID   string `header:"Trace-Id"`
	SpanID    string `header:"Span-Id"`

	Language language.Language `header:"Accept-Language"`
}
