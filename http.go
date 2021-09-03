package visiauth

import "net/http"

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type MocksHttpClient struct {
	fn func(url string) (resp *http.Response, err error)
}

func NewMocksHttpClient(fn func(url string) (resp *http.Response, err error)) *MocksHttpClient {
	return &MocksHttpClient{fn}
}

func (c *MocksHttpClient) Get(url string) (resp *http.Response, err error) {
	return c.fn(url)
}
