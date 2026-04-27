package errs

import (
	"fmt"
	"net/http"
)

type Err struct {
	Status      int    `json:"-"`
	ErrorCode   string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`

	Data map[string]any `json:"-"`
	Err  error          `json:"-"`
}

func (e *Err) Error() string {
	return fmt.Sprintf("status: %d, code: %s", e.Status, e.ErrorCode)
}

func (e *Err) HttpStatus() int {
	return e.Status
}

func (e *Err) Code() string {
	return e.ErrorCode
}

func New(status int, code string, errs ...error) Error {
	return &Errs{
		Errors: []Err{
			{
				Status:    status,
				ErrorCode: code,
			},
		},
	}
}

func NewCustom(status int, code, message, description string, errs ...error) Error {
	return &Errs{
		Errors: []Err{
			{
				Status:      status,
				ErrorCode:   code,
				Message:     message,
				Description: description,
				Err:         getError(errs...),
			},
		},
	}
}

func NewInternalError(errs ...error) Error {
	return &Errs{
		Errors: []Err{
			{
				Status:    http.StatusInternalServerError,
				ErrorCode: Err50001,
				Err:       getError(errs...),
			},
		},
	}
}

func NewExternalError(errs ...error) Error {
	return &Errs{
		Errors: []Err{
			{
				Status:    http.StatusServiceUnavailable,
				ErrorCode: Err50003,
				Err:       getError(errs...),
			},
		},
	}
}

func NewBadRequestError(errs ...error) Error {
	return &Errs{
		Errors: []Err{
			{
				Status:    http.StatusBadRequest,
				ErrorCode: Err40000,
				Err:       getError(errs...),
			},
		},
	}
}

func NewBusinessError(code string, errs ...error) Error {
	return &Errs{
		Errors: []Err{
			{
				Status:    http.StatusConflict,
				ErrorCode: code,
				Err:       getError(errs...),
			},
		},
	}
}

func getError(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return errs[0]
}
