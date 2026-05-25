package reqs

import (
	"net/http"

	"github.com/witwoywhy/req"
)

type Response interface {
	IsErrorState() bool
	Error() error
	HTTPStatus() int
	Header() http.Header
}

type response struct {
	response *req.Response
}

func (r response) Error() error {
	return r.response.Err
}

func (r response) IsErrorState() bool {
	return r.response.IsErrorState()
}

func (r response) HTTPStatus() int {
	return r.response.GetStatusCode()
}

func (r response) Header() http.Header {
	return r.response.Header
}
