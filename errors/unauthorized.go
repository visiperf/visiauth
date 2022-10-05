package errors

import (
	"errors"
	"net/http"
)

type unauthorized struct {
	reason, code string
}

func Unauthorized(reason, code string) error {
	return unauthorized{reason, code}
}

func (u unauthorized) Message() string {
	return u.reason
}

func (u unauthorized) Code() string {
	return u.code
}

func (u unauthorized) StatusCode() int {
	return http.StatusUnauthorized
}

func (u unauthorized) String() string {
	return u.Message()
}

func (u unauthorized) Error() string {
	return u.String()
}

func IsUnauthorized(err error) bool {
	var u unauthorized
	return errors.As(err, &u)
}
