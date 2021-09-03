package visiauth

import "github.com/golang-jwt/jwt"

type TokenParser interface {
	ParseToken(accessToken string) (*jwt.Token, error)
}

type JwtTokenParser struct {
	certificateFetcher PEMCertificateFetcher
}

func NewJwtTokenParser(certificateFetcher PEMCertificateFetcher) *JwtTokenParser {
	return &JwtTokenParser{certificateFetcher}
}

func (p *JwtTokenParser) ParseToken(accessToken string) (*jwt.Token, error) {
	return jwt.Parse(accessToken, p.keyFunc)
}

func (p *JwtTokenParser) keyFunc(token *jwt.Token) (interface{}, error) {
	cert, err := p.certificateFetcher.FetchPEMCertificate(token)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(cert)
}
