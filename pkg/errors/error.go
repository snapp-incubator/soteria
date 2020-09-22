package errors

import (
	"fmt"
)

// Error is the error type that contains a Code and a message about the error
type Error struct {
	Code    Code
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code.Message(), e.Message)
}

// CreateError returns an instance of Error with given information
func CreateError(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
