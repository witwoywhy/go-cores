package reqs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/utils"
	"github.com/witwoywhy/req"
)

type Request interface {
	SetContext(ctx context.Context) Request

	SetPathParam(key string, value string) Request
	SetPathParams(params map[string]string) Request
	AddQueryParam(key string, value string) Request
	AddQueryParams(key string, values ...string) Request
	SetHeader(key, value string) Request
	SetHeaders(hdrs map[string]string) Request
	SetBearerAuthToken(token string) Request
	SetBasicAuth(username, password string) Request
	SetFormData(data map[string]interface{}) Request
	SetFileReader(paramName, filename string, reader io.Reader) Request

	SetBody(body interface{}) Request
	SetError(err interface{}) Request
	SetResult(result interface{}) Request
	Do() Response
}

type request struct {
	request *req.Request
	config  *Config
	l       logger.Logger
}

func (r *request) SetContext(ctx context.Context) Request {
	r.request.SetContext(ctx)
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
	r.request.SetBodyJsonMarshal(body)
	return r
}

func (r *request) SetFormData(data map[string]interface{}) Request {
	r.request.SetFormDataAnyType(data)
	return r
}

func (r *request) SetFileReader(paramName, filename string, reader io.Reader) Request {
	r.request.SetFileReader(paramName, filename, reader)
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
	var (
		doResponse     *req.Response
		requestBody    map[string]any
		responseBody   map[string]any
		requestHeader  http.Header
		responseHeader http.Header
	)

	fullPath := fmt.Sprintf("%s%s", r.config.BaseUrl, r.config.Url)
	requestHeader = r.request.Headers.Clone()
	utils.MaskHeader(apps.Authorization, requestHeader)

	if len(r.request.Body) > 0 && !r.config.EnableIgnoreLogBody {
		err := json.Unmarshal(r.request.Body, &requestBody)
		if err != nil {
			r.l.Errorf("failed when json.Unmarshal request body: %v, %s", err, string(r.request.Body))
			return response{response: &req.Response{Err: err}}
		}
	}

	r.l.JSON(map[string]any{
		apps.Key:     apps.StartOutbound,
		apps.Header:  requestHeader,
		apps.Body:    requestBody,
		apps.Method:  r.config.Method,
		apps.URL:     fullPath,
		apps.Message: fmt.Sprintf(apps.StartOutboundFmt, r.config.Method, fullPath),
	})

	switch strings.ToUpper(r.config.Method) {
	case http.MethodGet:
		doResponse, _ = r.request.Get(r.config.Url)
	case http.MethodPost:
		doResponse, _ = r.request.Post(r.config.Url)
	case http.MethodPatch:
		doResponse, _ = r.request.Patch(r.config.Url)
	case http.MethodPut:
		doResponse, _ = r.request.Put(r.config.Url)
	case http.MethodDelete:
		doResponse, _ = r.request.Delete(r.config.Url)
	}

	var response response = response{
		response: doResponse,
	}

	if response.Error() != nil {
		return response
	}

	if doResponse == nil {
		return response
	}

	processTime := fmt.Sprintf("%v", doResponse.TotalTime())
	responseHeader = doResponse.Header.Clone()
	utils.MaskHeader(apps.Authorization, responseHeader)

	if !r.config.EnableIgnoreLogBody {
		b, err := doResponse.ToBytes()
		if err != nil {
			return response
		}

		if len(b) > 0 {
			err := json.Unmarshal(b, &responseBody)
			if err != nil {
				r.l.Errorf("failed to json.Unmarshal response body: %v", err)
			}
		}
	}

	r.l.JSON(map[string]any{
		apps.Key:         apps.EndOutbound,
		apps.Header:      responseHeader,
		apps.Body:        responseBody,
		apps.HTTPStatus:  doResponse.StatusCode,
		apps.ProcessTime: processTime,
		apps.URL:         r.request.URL.String(),
		apps.Message:     fmt.Sprintf(apps.EndOutboundFmt, doResponse.StatusCode, processTime, r.request.URL.String()),
	})
	r.l.JSON(map[string]any{
		apps.Key:            apps.SummaryOutbound,
		apps.Method:         r.config.Method,
		apps.RequestHeader:  requestHeader,
		apps.ResponseHeader: responseHeader,
		apps.RequestBody:    requestBody,
		apps.ResponseBody:   responseBody,
		apps.HTTPStatus:     doResponse.StatusCode,
		apps.ProcessTime:    processTime,
		apps.URL:            r.request.URL.String(),
		apps.Message:        fmt.Sprintf(apps.SummaryOutboundFmt, doResponse.StatusCode, processTime, r.request.URL.String()),
	})
	return response
}
