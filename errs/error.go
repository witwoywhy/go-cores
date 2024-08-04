package errs

import (
	"fmt"
	"net/http"
)

type Error interface {
	Error() string
	Status() int
	Code() string
}

type Errs struct {
	StatusCode  int    `json:"-"`
	ErrorCode   string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (e *Errs) Error() string {
	return fmt.Sprintf("status: %d, code: %s", e.Status, e.ErrorCode)
}

func (e *Errs) Status() int {
	return e.StatusCode
}

func (e *Errs) Code() string {
	return e.ErrorCode
}

func New(status int, code string) Error {
	return &Errs{
		StatusCode: status,
		ErrorCode:  code,
	}
}

func NewCustom(status int, code, message, description string) Error {
	return &Errs{
		StatusCode:  status,
		ErrorCode:   code,
		Message:     message,
		Description: description,
	}
}

func NewInternalError() Error {
	return &Errs{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  Err50001,
	}
}

func NewExternalError() Error {
	return &Errs{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  Err50001,
	}
}

func NewBadRequestError() Error {
	return &Errs{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  Err40000,
	}
}

func NewBusinessError(code string) Error {
	return &Errs{
		StatusCode: http.StatusConflict,
		ErrorCode:  code,
	}
}
