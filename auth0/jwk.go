package auth0

import (
	"context"
	"fmt"

	"github.com/visiperf/visiauth/v3"
)

type JwkFetcher struct {
	jwksFetcher *JwksFetcher
}

func NewJwkFetcher(domain string) *JwkFetcher {
	return &JwkFetcher{
		jwksFetcher: NewJwksFetcher(domain),
	}
}

func (f *JwkFetcher) FetchJwk(_ context.Context, kid string) (*visiauth.Jwk, error) {
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
