package visiauth

import (
	"context"
	"fmt"
)

type CertificateFetcher struct {
	jwkFetcher JwkFetcher
}

func NewCertificateFetcher(jwkFetcher JwkFetcher) *CertificateFetcher {
	return &CertificateFetcher{jwkFetcher}
}

func (f *CertificateFetcher) FetchPEMCertificate(ctx context.Context, token *UserToken) ([]byte, error) {
	jwk, err := f.jwkFetcher.FetchJwk(ctx, token.Kid())
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", jwk.X5c[0])), nil
}
