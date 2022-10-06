package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type notFound struct {
	resource, code string
}

func NotFound(resource, code string) error {
	return notFound{
		resource: resource,
		code:     code,
	}
}

func (nf notFound) Message() string {
	return fmt.Sprintf("%s not found", nf.resource)
}

func (nf notFound) Code() string {
	return nf.code
}

func (nf notFound) StatusCode() int {
	return http.StatusNotFound
}

func (nf notFound) String() string {
	return nf.Message()
}

func (nf notFound) Error() string {
	return nf.String()
}

func IsNotFound(err error) bool {
	var nf notFound
	return errors.As(err, &nf)
}
