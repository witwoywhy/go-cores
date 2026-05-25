package errs

type Error interface {
	Error() string
	HttpStatus() int
	Code() string
}
