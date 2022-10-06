package errors

import (
	"encoding/json"
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

func (u unauthorized) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}{
		Message: u.Message(),
		Code:    u.Code(),
	})
}

func IsUnauthorized(err error) bool {
	var u unauthorized
	return errors.As(err, &u)
}
