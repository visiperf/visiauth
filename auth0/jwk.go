package auth0

import (
	"context"
	"fmt"

	"github.com/visiperf/visiauth/v2"
)

type JwkFetcher struct {
	jwksFetcher *JwksFetcher
}

func NewJwkFetcher() *JwkFetcher {
	return &JwkFetcher{
		jwksFetcher: NewJwksFetcher(),
	}
}

func (f *JwkFetcher) FetchJwk(ctx context.Context, kid string) (*visiauth.Jwk, error) {
	jwks, err := f.jwksFetcher.FetchJwks(ctx)
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
