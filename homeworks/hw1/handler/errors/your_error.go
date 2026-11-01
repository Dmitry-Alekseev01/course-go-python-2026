package errors

import "fmt"

// AdditionalMessageError ошибка, которая
type AdditionalMessageError struct {
	err    error
	format string
	args   []any
}

// NewAdditionalMessageError создает ошибку
func NewAdditionalMessageError(err error, format string, args ...any) error {
	return &AdditionalMessageError{
		err:    err,
		format: format,
		args:   args,
	}
}

func (e *AdditionalMessageError) Error() string {
	message := fmt.Sprintf(e.format, e.args...)

	switch {
	case e.err == nil:
		return message
	case e.format == "":
		return e.err.Error()
	default:
		return message + ": " + e.err.Error()
	}
}

func (e *AdditionalMessageError) Unwrap() error {
	return e.err
}
