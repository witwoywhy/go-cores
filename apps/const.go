package apps

const (
	Message     = "message"
	MaskingChar = "*"

	StartInbound      = "START INBOUND"
	EndInbound        = "END INBOUND"
	SummaryInbound    = "SUMMARY INBOUND"
	StartInboundFmt   = StartInbound + " | %s | %s | %s"
	EndInboundFmt     = EndInbound + " | %v | %v | %s | %s"
	SummaryInboundFmt = SummaryInbound + " | %v | %v | %s | %s"

	StartOutbound      = "START OUTBOUND"
	EndOutbound        = "END OUTBOUND"
	SummaryOutbound    = "SUMMARY OUTBOUND"
	StartOutboundFmt   = StartOutbound + " | %s | %s"
	EndOutboundFmt     = EndOutbound + " | %d | %s | %s"
	SummaryOutboundFmt = SummaryOutbound + " | %d | %s | %s"

	Header         = "header"
	RequestHeader  = "request_header"
	ResponseHeader = "response_header"

	Body         = "body"
	RequestBody  = "request_body"
	ResponseBody = "response_body"

	Method      = "method"
	Host        = "host"
	URL         = "url"
	HTTPStatus  = "http_status"
	ProcessTime = "process_time"
	Key         = "key"

	Authorization = "Authorization"
	TraceID       = "trace_id"
	SpanID        = "span_id"

	Rctx     = "rctx"
	Logger   = "logger"
	Language = "language"
)

var HeaderMaskingList = map[string]bool{
	"authorization": true,
	"x-api-key":     true,
}
