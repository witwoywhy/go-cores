package gins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/utils"
	"go.opentelemetry.io/otel/attribute"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Log(ignoreLogBody map[string]bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()

		writer := &responseBodyWriter{ResponseWriter: ctx.Writer, body: &bytes.Buffer{}}
		ctx.Writer = writer

		request := ctx.Request
		l, span := NewRequestLog(ctx, request.URL.Path)
		ctx.Set("logger", l)

		var (
			isHasSpan       = span != nil
			isIgnoreLogBody = ignoreLogBody[ctx.FullPath()]

			requestBody    map[string]any
			responseBody   map[string]any
			requestHeader  http.Header
			responseHeader http.Header
		)

		requestHeader = request.Header.Clone()
		utils.MaskHeader(apps.Authorization, requestHeader)

		if isHasSpan {
			defer span.End()
		}

		b, _ := io.ReadAll(request.Body)
		if len(b) > 0 {
			request.Body = io.NopCloser(bytes.NewBuffer(b))
			if !isIgnoreLogBody {
				json.Unmarshal(b, &requestBody)
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

		ctx.Next()

		if isHasSpan {
			span.SetAttributes(attribute.Int("http.response.status_code", writer.Status()))
		}

		if len(writer.body.Bytes()) > 0 && !isIgnoreLogBody {
			json.Unmarshal(writer.body.Bytes(), &responseBody)
		}

		processTime := fmt.Sprintf("%v", time.Since(now))
		responseHeader = writer.Header().Clone()
		utils.MaskHeader(apps.Authorization, responseHeader)

		l.JSON(map[string]any{
			apps.Key:         apps.EndInbound,
			apps.Header:      responseHeader,
			apps.Body:        responseBody,
			apps.HTTPStatus:  writer.Status(),
			apps.ProcessTime: processTime,
			apps.Message:     fmt.Sprintf(apps.EndInboundFmt, writer.Status(), processTime, request.Method, request.URL.Path),
		})
		l.JSON(map[string]any{
			apps.Key:            apps.SummaryInbound,
			apps.Method:         request.Method,
			apps.RequestHeader:  requestHeader,
			apps.ResponseHeader: responseHeader,
			apps.RequestBody:    requestBody,
			apps.ResponseBody:   responseBody,
			apps.HTTPStatus:     writer.Status(),
			apps.ProcessTime:    processTime,
			apps.URL:            request.URL.Path,
			apps.Message:        fmt.Sprintf(apps.SummaryInboundFmt, writer.Status(), processTime, request.Method, request.URL.Path),
		})
	}
}
