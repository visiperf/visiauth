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

type JwksFetcher struct {
	domain string
}

func NewJwksFetcher(domain string) *JwksFetcher {
	return &JwksFetcher{domain}
}

func (f *JwksFetcher) FetchJwks() (*Jwks, error) {
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

func (f *JwksFetcher) url() string {
	return fmt.Sprintf(auth0JwksUrl, f.domain)
}
