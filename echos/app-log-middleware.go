package echos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/utils"
	"go.opentelemetry.io/otel/attribute"
)

type responseBodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseBodyWriter) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func Log(ignoreLogBody map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			now := time.Now()

			writer := &responseBodyWriter{ResponseWriter: ctx.Response(), body: &bytes.Buffer{}}
			ctx.SetResponse(writer)

			request := ctx.Request()
			l, span := NewRequestLog(ctx, request.URL.Path)
			ctx.Set("logger", l)

			var (
				isHasSpan       = span != nil
				isIgnoreLogBody = ignoreLogBody[request.URL.RawPath]

				requestBody    any
				responseBody   any
				requestHeader  http.Header
				responseHeader http.Header
			)

			if ctx.Path() != "" {
				isIgnoreLogBody = ignoreLogBody[ctx.Path()]
			}

			requestHeader = request.Header.Clone()
			utils.MaskHeader(apps.Authorization, requestHeader)

			if isHasSpan {
				defer span.End()
			}

			b, _ := io.ReadAll(request.Body)
			if len(b) > 0 {
				request.Body = io.NopCloser(bytes.NewBuffer(b))
				if !isIgnoreLogBody {
					requestBody = decodeBody(b)
				}
			}

			l.JSON(map[string]any{
				apps.Key:     apps.StartInbound,
				apps.Header:  requestHeader,
				apps.Body:    requestBody,
				apps.Method:  request.Method,
				apps.Host:    request.Host,
				apps.URL:     request.URL.Path,
				apps.Message: fmt.Sprintf(apps.StartInboundFmt, request.Method, request.Host, request.URL.Path),
			})

			err := next(ctx)
			statusCode := statusCode(ctx)

			if isHasSpan {
				span.SetAttributes(attribute.Int("http.response.status_code", statusCode))
			}

			if len(writer.body.Bytes()) > 0 && !isIgnoreLogBody {
				responseBody = decodeBody(writer.body.Bytes())
			}

			processTime := fmt.Sprintf("%v", time.Since(now))
			responseHeader = ctx.Response().Header().Clone()
			utils.MaskHeader(apps.Authorization, responseHeader)

			l.JSON(map[string]any{
				apps.Key:         apps.EndInbound,
				apps.Header:      responseHeader,
				apps.Body:        responseBody,
				apps.HTTPStatus:  statusCode,
				apps.ProcessTime: processTime,
				apps.Message:     fmt.Sprintf(apps.EndInboundFmt, statusCode, processTime, request.Method, request.URL.Path),
			})
			l.JSON(map[string]any{
				apps.Key:            apps.SummaryInbound,
				apps.Method:         request.Method,
				apps.RequestHeader:  requestHeader,
				apps.ResponseHeader: responseHeader,
				apps.RequestBody:    requestBody,
				apps.ResponseBody:   responseBody,
				apps.HTTPStatus:     statusCode,
				apps.ProcessTime:    processTime,
				apps.URL:            request.URL.Path,
				apps.Message:        fmt.Sprintf(apps.SummaryInboundFmt, statusCode, processTime, request.Method, request.URL.Path),
			})
			return err
		}
	}
}

func decodeBody(b []byte) any {
	var body any
	if err := json.Unmarshal(b, &body); err == nil {
		return body
	}

	return string(b)
}

func statusCode(ctx *echo.Context) int {
	response, err := echo.UnwrapResponse(ctx.Response())
	if err != nil {
		return http.StatusOK
	}

	if response.Status == 0 {
		return http.StatusOK
	}

	return response.Status
}
