package contexts

import "github.com/witwoywhy/go-cores/enum/language"

type Header struct {
	Authorization string `header:"Authorization"`

	RequestID string `header:"X-Request-Id"`
	ClientID  string `header:"X-Client-Id"`

	TraceID string `header:"Trace-Id"`
	SpanID  string `header:"Span-Id"`

	UserAgent string            `header:"User-Agent"`
	Language  language.Language `header:"Accept-Language"`
}
