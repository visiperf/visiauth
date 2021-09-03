package visiauth

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestWrapChainToCertificateConverter(t *testing.T) {
	tests := []struct {
		name        string
		converter   *WrapChainToCertificateConverter
		chain       string
		certificate string
	}{{
		name:      "success",
		converter: NewWrapChainToCertificateConverter(),
		chain:     `azerty`,
		certificate: `-----BEGIN CERTIFICATE-----
		azerty
		-----END CERTIFICATE-----`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, strings.ReplaceAll(test.certificate, "\t", ""), test.converter.ConvertChainToCertificate(test.chain))
		})
	}
}

func TestAuth0PEMCertificateFetcher(t *testing.T) {
	tests := []struct {
		name    string
		fetcher *Auth0PEMCertificateFetcher
		token   *jwt.Token
		cert    []byte
		err     error
	}{{
		name: "connection unavailable",
		fetcher: NewAuth0PEMCertificateFetcher("domain.io", NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return nil, errors.New("connection unavailable")
		})),
		token: &jwt.Token{Header: map[string]interface{}{
			"kid": "azerty",
		}},
		cert: nil,
		err:  errors.New("connection unavailable"),
	}, {
		name: "success",
		fetcher: NewAuth0PEMCertificateFetcher("domain.io", NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(jwksJsonMocks))),
			}, nil
		})),
		token: &jwt.Token{Header: map[string]interface{}{
			"kid": "rF45rmcRt-gEXpMBzvw3U",
		}},
		cert: []byte("-----BEGIN CERTIFICATE-----\n" + jwksMocks.Keys[0].X5c[0] + "\n-----END CERTIFICATE-----"),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cert, err := test.fetcher.FetchPEMCertificate(test.token)

			assert.Equal(t, test.cert, cert)
			assert.Equal(t, test.err, err)
		})
	}
}
