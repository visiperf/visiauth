package visiauth

import "net/http"

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type MocksHttpClient struct {
	get func(url string) (resp *http.Response, err error)
}

func NewMocksHttpClient(get func(url string) (resp *http.Response, err error)) *MocksHttpClient {
	return &MocksHttpClient{get}
}

func (c *MocksHttpClient) Get(url string) (resp *http.Response, err error) {
	return c.get(url)
}
