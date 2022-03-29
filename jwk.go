package visiauth

import (
	"context"
	"encoding/json"
	"fmt"
)

type Jwk struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func (jwk Jwk) Is(kid string) bool {
	return jwk.Kid == kid
}

func (jwk Jwk) MarshalBinary() ([]byte, error) {
	return json.Marshal(jwk)
}

func (jwk *Jwk) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, jwk)
}

type JwkFetcher interface {
	FetchJwk(ctx context.Context, kid string) (*Jwk, error)
}

type Auth0JwkFetcher struct {
	jwksFetcher *JwksFetcher
}

func NewAuth0JwkFetcher(domain string) *Auth0JwkFetcher {
	return &Auth0JwkFetcher{
		jwksFetcher: NewJwksFetcher(domain),
	}
}

func (f *Auth0JwkFetcher) FetchJwk(_ context.Context, kid string) (*Jwk, error) {
	jwks, err := f.jwksFetcher.FetchJwks()
	if err != nil {
		return nil, err
	}

	for _, jwk := range jwks.Keys {
		if jwk.Is(kid) {
			return jwk, nil
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
