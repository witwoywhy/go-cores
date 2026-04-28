package reqs

import "github.com/witwoywhy/req"

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
