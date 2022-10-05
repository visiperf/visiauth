package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type invalidArgument struct {
	reason, code string
}

func InvalidArgument(reason, code string) error {
	return invalidArgument{reason, code}
}

func (ia invalidArgument) Message() string {
	return ia.reason
}

func (ia invalidArgument) Code() string {
	return ia.code
}

func (ia invalidArgument) StatusCode() int {
	return http.StatusBadRequest
}

func (ia invalidArgument) String() string {
	return ia.Message()
}

func (ia invalidArgument) Error() string {
	return ia.String()
}

func (ia invalidArgument) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}{
		Message: ia.Message(),
		Code:    ia.Code(),
	})
}

func IsInvalidArgument(err error) bool {
	var ia invalidArgument
	return errors.As(err, &ia)
}
