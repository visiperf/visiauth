package visiauth

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonDecoder(t *testing.T) {
	type obj struct {
		Key string `json:"key"`
	}

	tests := []struct {
		name    string
		decoder *JsonDecoder
		reader  io.ReadCloser
		res     interface{}
		err     error
	}{{
		name:    "invalid json",
		decoder: NewJsonDecoder(),
		reader:  io.NopCloser(bytes.NewBuffer([]byte(`{ "key": value }`))),
		res:     obj{},
		err:     &json.SyntaxError{},
	}, {
		name:    "success",
		decoder: NewJsonDecoder(),
		reader:  io.NopCloser(bytes.NewBuffer([]byte(`{ "key": "value" }`))),
		res:     obj{Key: "value"},
		err:     nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o obj
			err := test.decoder.Decode(test.reader, &o)

			assert.Equal(t, test.res, o)

			if test.err != nil {
				assert.ErrorAs(t, err, &test.err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
