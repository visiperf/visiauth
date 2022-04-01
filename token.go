package visiauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Token struct {
	*jwt.Token
}

func NewToken(token *jwt.Token) *Token {
	return &Token{token}
}

func (t Token) Header() map[string]interface{} {
	return t.Token.Header
}

func (t Token) Kid() string {
	return t.Header()["kid"].(string)
}

func (t Token) Claims() map[string]interface{} {
	return t.Token.Claims.(jwt.MapClaims)
}

func (t Token) UserID() string {
	return t.customKey("user_id")
}

func (t Token) Iss() string {
	return t.Claims()["iss"].(string)
}

func (t Token) Scopes() []string {
	return strings.Split(t.scope(), " ")
}

func (t Token) scope() string {
	return t.Claims()["scope"].(string)
}

func (t Token) customKey(key string) string {
	return fmt.Sprintf("https://%s%s", t.Iss(), key)
}

type TokenParser struct {
	certificateFetcher *CertificateFetcher
}

func NewTokenParser(jwkFetcher JwkFetcher) *TokenParser {
	return &TokenParser{NewCertificateFetcher(jwkFetcher)}
}

func (p *TokenParser) ParseToken(ctx context.Context, accessToken string) (*jwt.Token, error) {
	return jwt.Parse(accessToken, p.keyFunc(ctx))
}

func (p *TokenParser) keyFunc(ctx context.Context) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		cert, err := p.certificateFetcher.FetchPEMCertificate(ctx, token)
		if err != nil {
			return nil, err
		}

		return jwt.ParseRSAPublicKeyFromPEM(cert)
	}
}
