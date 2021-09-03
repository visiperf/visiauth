package visiauth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type PEMCertificateFetcher interface {
	FetchPEMCertificate(token *jwt.Token) ([]byte, error)
}

type Auth0PEMCertificateFetcher struct {
	domain string
}

func NewAuth0PEMCertificateFetcher(domain string) *Auth0PEMCertificateFetcher {
	return &Auth0PEMCertificateFetcher{domain}
}

func (f *Auth0PEMCertificateFetcher) FetchPEMCertificate(token *jwt.Token) ([]byte, error) {
	panic("not implemented") // TODO: Implement
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
