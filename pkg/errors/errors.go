package errors

import (
	"fmt"
)

type CustomError struct {
	Code    int
	Message string
}

func NewError(message string, code int) error {
	return &CustomError{Message: message, Code: code}
}

func Error2Custom(err error) CustomError {
	customError, ok := err.(*CustomError)
	if !ok {
		return CustomError{Code: 500, Message: "invalid error message"}
	}
	return *customError
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
