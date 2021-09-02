package visiauth

import (
	"errors"
	"fmt"
)

type Jwks struct {
	Keys []Jwk `json:"keys"`
}

type Jwk struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type JwksFetcher interface {
	FetchJwks() (*Jwks, error)
}

type Auth0JwksFetcher struct {
	domain  string
	client  HttpClient
	decoder Decoder
}

func NewAuth0JwksFetcher(domain string, client HttpClient) *Auth0JwksFetcher {
	return &Auth0JwksFetcher{
		domain:  domain,
		client:  client,
		decoder: NewJsonDecoder(),
	}
}

func (f *Auth0JwksFetcher) FetchJWKS() (*Jwks, error) {
	resp, err := f.client.Get(f.url())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}

	var jwks Jwks
	if err := f.decoder.Decode(resp.Body, &jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func (f *Auth0JwksFetcher) url() string {
	return fmt.Sprintf("https://%s/.well-known/jwks.json", f.domain)
}
