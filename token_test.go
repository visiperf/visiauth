package visiauth

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestRetrieveTokenFromContext(t *testing.T) {
	tests := []struct {
		name  string
		ctx   context.Context
		token string
		err   error
	}{{
		name: "empty metadata",
		ctx:  context.Background(),
		err:  ErrMissingMetadata,
	}, {
		name: "empty authorization",
		ctx:  metadata.NewIncomingContext(context.Background(), metadata.MD{}),
		err:  ErrMissingAuthorization,
	}, {
		name:  "ok",
		ctx:   metadata.NewIncomingContext(context.Background(), metadata.Pairs(AuthorizationKey, "Bearer abc.def.ghi")),
		token: "abc.def.ghi",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := RetrieveTokenFromContext(test.ctx)

			assert.Equal(t, test.token, token)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestRetrieveTokenFromRequest(t *testing.T) {
	tests := []struct {
		name  string
		req   *http.Request
		token string
		err   error
	}{{
		name: "empty authorization",
		req:  &http.Request{},
		err:  ErrMissingAuthorization,
	}, {
		name: "ok",
		req: &http.Request{
			Header: http.Header{
				AuthorizationKey: {"Bearer abc.def.ghi"},
			},
		},
		token: "abc.def.ghi",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := RetrieveTokenFromRequest(test.req)

			assert.Equal(t, test.token, token)
			assert.Equal(t, test.err, err)
		})
	}
}
