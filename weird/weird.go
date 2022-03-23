// Package weird is our internal error handling tool
// with proper protocol agnostic response definition
package weird

import "errors"

// Error is the main error handling representation, the usage
type Error struct {
	// StatusCode is the final status code returned for http response writer
	StatusCode int

	// Error is the original error
	InnerError error

	// Msg is the final message returned for http response writer
	Msg string
}

func New(msg string, err error, statusCode int) error {
	return Error{
		StatusCode: statusCode,
		InnerError: err,
		Msg:        msg,
	}
}

// Error satisfy the error interface and format the error
func (e Error) Error() string {
	return e.InnerError.Error()
}

// Is helps to manage the errors.Is
func (e Error) Is(target error) bool {
	return errors.Is(e.InnerError, target)
}
