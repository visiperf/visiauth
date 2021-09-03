package visiauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jwksMocks Jwks

func init() {
	if err := json.Unmarshal([]byte(jwksJsonMocks), &jwksMocks); err != nil {
		log.Fatal(err)
	}
}

const jwksJsonMocks = `
{
	"keys": [{
		"alg": "RS256",
		"kty": "RSA",
		"use": "sig",
		"n": "1uLJYzaC4k-8HA16AblDoaB4g8w6qvSyY6Ai43EuoPVUKTjTtij8V9yRfV_F2xVXBUmx26UK18w-LcIYGxAf33Zxs1K9GZ_SXDgui1Yy9y7iYk_Tm57pmQ_pxT6ZY0p5mGVxQHTsF6HpBmW7GabUF-PDJ8nxEjzfl72JPegzZ1xWnv_RTGn8ohFfvEiAH95EAe3gH9xwWvP3fFEfxnq0ouVUgZKzr_P_BYs-UM4ZmJgGLth2PPiyhmLtT0kCEo5zb4Xv4eFDAYbJmyCV9yBAaNapPr1idJ5rRaQpt2yCi5AlQ8xmU5zpcmMQFF8Vy78JSgH6KBdKgpdOm6Mv6RJZvQ",
		"e": "AQAB",
		"kid": "rF45rmcRt-gEXpMBzvw3U",
		"x5t": "wM4fWABeXT0WwC99cDH5gbebp2o",
		"x5c": ["MIIDDTCCAfWgAwIBAgIJfmQjGqVWSTThMA0GCSqGSIb3DQEBCwUAMCQxIjAgBgNVBAMTGWRldi12aXNpcGVyZi5ldS5hdXRoMC5jb20wHhcNMjEwODA1MDkyNTIxWhcNMzUwNDE0MDkyNTIxWjAkMSIwIAYDVQQDExlkZXYtdmlzaXBlcmYuZXUuYXV0aDAuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1uLJYzaC4k+8HA16AblDoaB4g8w6qvSyY6Ai43EuoPVUKTjTtij8V9yRfV/F2xVXBUmx26UK18w+LcIYGxAf33Zxs1K9GZ/SXDgui1Yy9y7iYk/Tm57pmQ/pxT6ZY0p5mGVxQHTsF6HpBmW7GabUF+PDJ8nxEjzfl72JPegzZ1xWnv/RTGn8ohFfvEiAH95EAe3gH9xwWvP3fFEfxnq0ouVUgZKzr/P/BYs+UM4ZmJgGLth2PPiyhmLtT0kCEo5zb4Xv4eFDAYbJmyCV9yBAaNapPr1idJ5rRaQpt2yCi5AlQ8xmU5zpcmMQFF8Vy78JSgH6KBdKgpdOm6Mv6RJZvQIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBR/C19XBb5LwoHecyxDwkPA10wSmjAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBADN9BpPkDAU/jAitgNkJdzrdOoNgVOOvCxlvB+SQKz2Pc2V4QRY4ZBrTdS+ig5/c0FBYWeKBwyXrxadoflYGCi4iye2jHJE5waHRoZ75T+bydUC1a189316W+xTpUL5STi9b5eJ44D+tuSnL+DIIRHZtFQEf5KbgcdZVbI2m2dbJJKnnkGv0VwPj2g1W6OFl1gUIVve4diMc5j/j4hWxkyB8AHCuu5WPIrYKHqmFn5qvalS0T8BFlD8Mhpg0cGdlEERxvN0XKON/sk+lcoid4sYqTqFvCsRprdPAprzCCyzRWxyCyPcuTDPTRfhHwnPzi+5DCRlespv/xnXNAsmXjMU="]
	}, {
		"alg": "RS256",
		"kty": "RSA",
		"use": "sig",
		"n": "rEWg2BtXa3TOAtJMw_4dezSbC7ITc2RnOtHLJVc2sKqM1eIBO5WsrdPaM0gXt9Ru9J_2CRvi-kiIk2OhPNr8iI_uHAh_aZbMUl5YChcAruc3xGwH9lOfwR3QwcwMi3USsRSX1VCpa9epLThd6tebdpX-IFUIYyxskzLpaSLPppjnMt-dsx9rASkAwGXFlO8FTB4ZKqj6HL_Z7I_mDRcD19fAA9oESOB6k8mkAfqOE6SFtKLiledtWEgSrHQqCjrOxawF0A71zQO4LMbo7TDoglJt4MEoWgaCgAOYK7CRR1Xsyxtq0KGWTxEut5j82H8inBzdxd8Ie3WV4G7-ccwHuw",
		"e": "AQAB",
		"kid": "ZDGPiuW8UNp_NZXlAUtyK",
		"x5t": "-d1M0AeQEmW1bQUP5vZHVuSKfj0",
		"x5c": ["MIIDDTCCAfWgAwIBAgIJQMq/3JQ6QL15MA0GCSqGSIb3DQEBCwUAMCQxIjAgBgNVBAMTGWRldi12aXNpcGVyZi5ldS5hdXRoMC5jb20wHhcNMjEwODA1MDkyNTIxWhcNMzUwNDE0MDkyNTIxWjAkMSIwIAYDVQQDExlkZXYtdmlzaXBlcmYuZXUuYXV0aDAuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArEWg2BtXa3TOAtJMw/4dezSbC7ITc2RnOtHLJVc2sKqM1eIBO5WsrdPaM0gXt9Ru9J/2CRvi+kiIk2OhPNr8iI/uHAh/aZbMUl5YChcAruc3xGwH9lOfwR3QwcwMi3USsRSX1VCpa9epLThd6tebdpX+IFUIYyxskzLpaSLPppjnMt+dsx9rASkAwGXFlO8FTB4ZKqj6HL/Z7I/mDRcD19fAA9oESOB6k8mkAfqOE6SFtKLiledtWEgSrHQqCjrOxawF0A71zQO4LMbo7TDoglJt4MEoWgaCgAOYK7CRR1Xsyxtq0KGWTxEut5j82H8inBzdxd8Ie3WV4G7+ccwHuwIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBS4Z9YAskfwn0+f3KG85qp5IDqEaDAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBAFNZT8i2Ij2OqZlDb4bmAGUe3Y35dKMWXF19UVrmr1nYmCTU4tAPCilhf35FspKzlLWTbSL5WmUl8DQQv4TY4d3mBqFDKckz3/RdHcrb8yDo+Cxfldb8GCbHJnmxwwTCJqADaHDzAmB/PlZAwqPYzNtj+M2cCJAv2+jncW9wHiuM+AowQNloV61eVBtLDJ8q4sIVxKjIjWgVRSjLr4hDsvgLKzrxDhpdjzKCQ1CMcSWiL9sr/L3B0aT+cB8zEqnIFon97Y5qBYMzXWHmQK1niHnxXlubCEP2QIdptyXss6O5DfWr1p9+Kh8+wAhb/qZS9pP8IcXVaPYrhrJFX3rQt7M="]
	}]
}
`

func TestAuth0JwksFetcher(t *testing.T) {
	tests := []struct {
		name    string
		fetcher *Auth0JwksFetcher
		jwks    *Jwks
		err     error
	}{{
		name: "connection unavailable",
		fetcher: NewAuth0JwksFetcher("domain.io", NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return nil, errors.New("connection unavailable")
		})),
		jwks: nil,
		err:  errors.New("connection unavailable"),
	}, {
		name: "not found",
		fetcher: NewAuth0JwksFetcher("domain.io", NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return &http.Response{
				Status:     "404 Not Found",
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
			}, nil
		})),
		jwks: nil,
		err:  errors.New("404 Not Found"),
	}, {
		name: "success",
		fetcher: NewAuth0JwksFetcher("domain.io", NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(jwksJsonMocks))),
			}, nil
		})),
		jwks: &jwksMocks,
		err:  nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jwks, err := test.fetcher.FetchJwks()

			assert.Equal(t, test.jwks, jwks)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestAuth0JwkFetcher(t *testing.T) {
	tests := []struct {
		name   string
		client HttpClient
		kid    string
		jwk    *Jwk
		err    error
	}{{
		name: "error with client",
		client: NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return nil, errors.New("connection unavailable")
		}),
		kid: "azerty",
		jwk: nil,
		err: errors.New("connection unavailable"),
	}, {
		name: "kid not found",
		client: NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(jwksJsonMocks))),
			}, nil
		}),
		kid: "azerty",
		jwk: nil,
		err: NewJwkNotFoundError("azerty"),
	}, {
		name: "success",
		client: NewMocksHttpClient(func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(jwksJsonMocks))),
			}, nil
		}),
		kid: "rF45rmcRt-gEXpMBzvw3U",
		jwk: &jwksMocks.Keys[0],
		err: nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jwk, err := NewAuth0JwkFetcher("domain.io", test.client).FetchJwk(test.kid)

			assert.Equal(t, test.jwk, jwk)
			assert.Equal(t, test.err, err)
		})
	}
}
