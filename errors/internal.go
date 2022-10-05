package errors

import (
	"errors"
	"net/http"
)

type internal struct {
	err error
}

func Internal(err error) error {
	return internal{err}
}

func (i internal) Message() string {
	return i.err.Error()
}

func (i internal) StatusCode() int {
	return http.StatusInternalServerError
}

func (i internal) String() string {
	return i.Message()
}

func (i internal) Error() string {
	return i.String()
}

func IsInternal(err error) bool {
	var i internal
	return errors.As(err, &i)
}
