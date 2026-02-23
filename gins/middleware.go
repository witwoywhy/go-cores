package gins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/logs"
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

func Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()

		writer := &responseBodyWriter{ResponseWriter: ctx.Writer, body: &bytes.Buffer{}}
		ctx.Writer = writer

		request := ctx.Request
		cloneHandler := request.Header.Clone()
		apps.MaskHeader(cloneHandler)

		l, span := NewRequestLog(ctx, request.URL.Path)
		isHasSpan := span != nil
		
		if isHasSpan {
			defer span.End()
		}

		ctx.Set("logger", l)

		requestBody := map[string]any{}
		b, _ := io.ReadAll(request.Body)
		if len(b) > 0 {
			request.Body = io.NopCloser(bytes.NewBuffer(b))
			json.Unmarshal(b, &requestBody)
		}

		l.JSON(map[string]any{
			apps.Header:  cloneHandler,
			apps.Body:    requestBody,
			logs.Message: fmt.Sprintf(apps.StartInbound, request.Method, request.Host, request.URL.Path),
		})

		ctx.Next()

		if isHasSpan {
			span.SetAttributes(attribute.Int("http.response.status_code", writer.Status()))
		}

		responseBody := map[string]any{}
		if len(writer.body.Bytes()) > 0 {
			json.Unmarshal(writer.body.Bytes(), &responseBody)
		}

		l.JSON(map[string]any{
			apps.Header:  writer.Header(),
			apps.Body:    responseBody,
			logs.Message: fmt.Sprintf(apps.EndInbound, writer.Status(), time.Since(now), request.Method, request.URL.Path),
		})
	}
}

func MakeDefaultRouteContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(rctx, &contexts.RouteContext{})
		ctx.Next()
	}
}
