package contexts

type Header struct {
	Authorization string
	
	TraceId       string
	SpanId        string
}
