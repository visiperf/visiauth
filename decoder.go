package visiauth

import (
	"encoding/json"
	"io"
)

type Decoder interface {
	Decode(reader io.ReadCloser, v interface{}) error
}

type JsonDecoder struct{}

func NewJsonDecoder() *JsonDecoder {
	return &JsonDecoder{}
}

func (d *JsonDecoder) Decode(reader io.ReadCloser, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}
