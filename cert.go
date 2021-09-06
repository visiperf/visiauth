package visiauth

import (
	"fmt"
)

type PEMCertificateFetcher interface {
	FetchPEMCertificate(token Token) ([]byte, error)
}

type Auth0PEMCertificateFetcher struct {
	fetcher   JwkFetcher
	converter ChainToCertificateConverter
}

func NewAuth0PEMCertificateFetcher(domain string, client HttpClient) *Auth0PEMCertificateFetcher {
	return &Auth0PEMCertificateFetcher{
		fetcher:   NewAuth0JwkFetcher(domain, client),
		converter: NewWrapChainToCertificateConverter(),
	}
}

func (f *Auth0PEMCertificateFetcher) FetchPEMCertificate(token Token) ([]byte, error) {
	jwk, err := f.fetcher.FetchJwk(token.Header().Kid())
	if err != nil {
		return nil, err
	}

	return []byte(f.converter.ConvertChainToCertificate(jwk.X5c[0])), nil
}

type ChainToCertificateConverter interface {
	ConvertChainToCertificate(chain string) string
}

type WrapChainToCertificateConverter struct{}

func NewWrapChainToCertificateConverter() *WrapChainToCertificateConverter {
	return &WrapChainToCertificateConverter{}
}

func (c *WrapChainToCertificateConverter) ConvertChainToCertificate(chain string) string {
	return fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", chain)
}
