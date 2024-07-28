package apps

const (
	TraceID = "traceId"
	SpanID  = "spanId"
)

const (
	StartInbound = "START INBOUND | %s | %s | %s"
	EndInbound   = "END INBOUND | %v | %v | %s | %s"

	StartOutbound = "START OUTBOUND | %s | %s"
	EndOutbound   = "END OUTBOUND | %d | %s | %s"
)

const (
	Header = "header"
	Body   = "body"
)

var HeaderMaskingList = map[string]bool{
	"Authorization": true,
}
