package auth0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/visiperf/visiauth/v3"
)

const (
	jwksUrl = "https://%s/.well-known/jwks.json"
)

type JwksFetcher struct {
	domain string
}

func NewJwksFetcher() *JwksFetcher {
	return &JwksFetcher{env.Auth0.Domain}
}

func (f *JwksFetcher) FetchJwks(_ context.Context) (*visiauth.Jwks, error) {
	resp, err := http.Get(fmt.Sprintf(jwksUrl, f.domain))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if isStatusCodeError(resp.StatusCode) {
		return nil, errors.New(resp.Status)
	}

	var jwks visiauth.Jwks
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}
