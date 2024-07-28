package reqs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
)

type Request interface {
	AddLogger(logger.Logger) Request
	SetPathParam(key string, value string) Request
	SetPathParams(params map[string]string) Request
	AddQueryParam(key string, value string) Request
	AddQueryParams(key string, values ...string) Request
	SetHeader(key, value string) Request
	SetHeaders(hdrs map[string]string) Request
	SetBearerAuthToken(token string) Request
	SetBasicAuth(username, password string) Request
	SetError(err interface{}) Request
	SetBody(body interface{}) Request
	SetResult(result interface{}) Request
	Do() Response
}

type request struct {
	request *req.Request
	config  *Config
	l       logger.CoreLogger
}

func (r *request) AddLogger(l logger.Logger) Request {
	log := l.(*logs.Log)
	r.l = logs.NewCoreLog(log.Information)
	return r
}

func (r *request) SetPathParam(key string, value string) Request {
	r.request.SetPathParam(key, value)
	return r
}

func (r *request) SetPathParams(params map[string]string) Request {
	r.request.SetPathParams(params)
	return r
}

func (r *request) AddQueryParam(key string, value string) Request {
	r.request.AddQueryParam(key, value)
	return r
}

func (r *request) AddQueryParams(key string, values ...string) Request {
	r.request.AddQueryParams(key, values...)
	return r
}

func (r *request) SetBody(body interface{}) Request {
	r.request.SetBody(body)
	return r
}

func (r *request) SetHeader(key string, value string) Request {
	r.request.SetHeader(key, value)
	return r
}

func (r *request) SetHeaders(hdrs map[string]string) Request {
	r.request.SetHeaders(hdrs)
	return r
}

func (r *request) SetBearerAuthToken(token string) Request {
	r.request.SetBearerAuthToken(token)
	return r
}

func (r *request) SetBasicAuth(username string, password string) Request {
	r.request.SetBasicAuth(username, password)
	return r
}

func (r *request) SetError(err interface{}) Request {
	r.request.SetErrorResult(err)

	return r
}

func (r *request) SetResult(result interface{}) Request {
	r.request.SetSuccessResult(result)
	return r
}

func (r *request) Do() Response {
	var doResponse *req.Response

	if r.l != nil {
		mapRequestBody := map[string]any{
			logs.Message: fmt.Sprintf(apps.StartOutbound, r.config.Api.Method, fmt.Sprintf("%s%s", r.config.BaseUrl, r.config.Api.Url)),
		}

		cloneRequestHeader := r.request.Headers.Clone()
		apps.MaskHeader(cloneRequestHeader)
		mapRequestBody[apps.Header] = cloneRequestHeader

		if len(r.request.Body) > 0 {
			var requestBody map[string]any
			json.Unmarshal(r.request.Body, &requestBody)
			mapRequestBody[apps.Body] = requestBody
		}

		r.l.JSON(mapRequestBody)
	}

	switch strings.ToLower(r.config.Api.Method) {
	case strings.ToLower(http.MethodGet):
		doResponse, _ = r.request.Get(r.config.Api.Url)
	case strings.ToLower(http.MethodPost):
		doResponse, _ = r.request.Post(r.config.Api.Url)
	case strings.ToLower(http.MethodPatch):
		doResponse, _ = r.request.Patch(r.config.Api.Url)
	case strings.ToLower(http.MethodPut):
		doResponse, _ = r.request.Put(r.config.Api.Url)
	case strings.ToLower(http.MethodDelete):
		doResponse, _ = r.request.Delete(r.config.Api.Url)
	}

	response := response{response: doResponse}
	if response.Error() != nil {
		return response
	}

	if r.l != nil {
		mapResponseBody := map[string]any{
			logs.Message: fmt.Sprintf(apps.EndOutbound, doResponse.StatusCode, doResponse.TotalTime(), r.request.URL.String()),
		}

		cloneResponseHeader := doResponse.Header.Clone()
		apps.MaskHeader(cloneResponseHeader)
		mapResponseBody[apps.Header] = cloneResponseHeader

		b, err := doResponse.ToBytes()
		if err != nil {
			return response
		}

		if len(b) > 0 {
			var responseBody map[string]any
			err := json.Unmarshal(b, &responseBody)
			if err != nil {
				r.l.Errorf("failed to unmarshal response body: %v", err)
			}

			if len(responseBody) > 0 {
				mapResponseBody[apps.Body] = responseBody
			}
		}

		r.l.JSON(mapResponseBody)
	}

	return response
}
