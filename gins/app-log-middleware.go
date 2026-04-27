package gins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/go-core-library.git/apps"
	"github.com/gin-gonic/gin"
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

func (a *app) Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()

		writer := &responseBodyWriter{ResponseWriter: ctx.Writer, body: &bytes.Buffer{}}
		ctx.Writer = writer

		request := ctx.Request
		l, span := NewRequestLog(ctx, request.URL.Path)
		ctx.Set("logger", l)

		var (
			isHasSpan       = span != nil
			isIgnoreLogBody = a.ignoreLogBody[ctx.FullPath()]

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
			apps.Url:     request.URL.Path,
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
			apps.HttpStatus:  writer.Status(),
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
			apps.HttpStatus:     writer.Status(),
			apps.ProcessTime:    processTime,
			apps.Url:            request.URL.Path,
			apps.Message:        fmt.Sprintf(apps.SummaryInboundFmt, writer.Status(), processTime, request.Method, request.URL.Path),
		})
	}
}
