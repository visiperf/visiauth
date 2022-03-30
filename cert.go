package visiauth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type CertificateFetcher struct {
	jwkFetcher JwkFetcher
}

func NewCertificateFetcher(jwkFetcher JwkFetcher) *CertificateFetcher {
	return &CertificateFetcher{jwkFetcher}
}

func (f *CertificateFetcher) FetchPEMCertificate(ctx context.Context, token jwt.Token) ([]byte, error) {
	jwk, err := f.jwkFetcher.FetchJwk(ctx, token.Header["kid"].(string))
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", jwk.X5c[0])), nil
}
