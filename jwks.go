package visiauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	auth0JwksUrl = "https://%s/.well-known/jwks.json"
)

type Jwks struct {
	Keys []*Jwk `json:"keys"`
}

type JwksFetcher interface {
	FetchJwks() (*Jwks, error)
}

type Auth0JwksFetcher struct {
	domain string
}

func NewAuth0JwksFetcher(domain string) *Auth0JwksFetcher {
	return &Auth0JwksFetcher{domain}
}

func (f *Auth0JwksFetcher) FetchJwks() (*Jwks, error) {
	resp, err := http.Get(f.url())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if isStatusCodeError(resp.StatusCode) {
		return nil, errors.New(resp.Status)
	}

	var jwks Jwks
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func (f *Auth0JwksFetcher) url() string {
	return fmt.Sprintf(auth0JwksUrl, f.domain)
}
