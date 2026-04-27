package errs

import "fmt"

type Errs struct {
	Errors []Err `json:"errors"`
}

func (e *Errs) Error() string {
	return fmt.Sprintf("status: %d, code: %s", e.Errors[0].Status, e.Errors[0].ErrorCode)
}

func (e *Errs) HttpStatus() int {
	return e.Errors[0].Status
}

func (e *Errs) Code() string {
	return e.Errors[0].ErrorCode
}

func (e *Errs) Append(err Err) {
	e.Errors = append(e.Errors, err)
}
