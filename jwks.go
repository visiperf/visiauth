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

func (f *Auth0JwksFetcher) FetchJwks() (*Jwks, error) {
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

type JwkFetcher interface {
	FetchJwk(kid string) (*Jwk, error)
}

type Auth0JwkFetcher struct {
	domain  string
	fetcher JwksFetcher
}

func NewAuth0JwkFetcher(domain string, client HttpClient) *Auth0JwkFetcher {
	return &Auth0JwkFetcher{
		domain:  domain,
		fetcher: NewAuth0JwksFetcher(domain, client),
	}
}

func (f *Auth0JwkFetcher) FetchJwk(kid string) (*Jwk, error) {
	jwks, err := f.fetcher.FetchJwks()
	if err != nil {
		return nil, err
	}

	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			return &jwk, nil
		}
	}

	return nil, NewJwkNotFoundError(kid)
}

type JwkNotFoundError struct {
	kid string
}

func NewJwkNotFoundError(kid string) *JwkNotFoundError {
	return &JwkNotFoundError{kid}
}

func (e *JwkNotFoundError) Error() string {
	return fmt.Sprintf("JWK with kid %s cannot be found", e.kid)
}
