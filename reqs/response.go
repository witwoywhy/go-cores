package reqs

import "github.com/imroc/req/v3"

type Response interface {
	IsErrorState() bool
	Error() error
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
