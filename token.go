package visiauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationKey    = "Authorization"
	authorizationPrefix = "Bearer "
)

var (
	ErrMissingMetadata      = errors.New("missing metadata")
	ErrMissingAuthorization = errors.New("missing authorization")
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
	return t.Claims()[t.customKey("user_id")].(string)
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
	return fmt.Sprintf("%s%s", t.Iss(), key)
}

type TokenParser struct {
	certificateFetcher *CertificateFetcher
}

func NewTokenParser(jwkFetcher JwkFetcher) *TokenParser {
	return &TokenParser{NewCertificateFetcher(jwkFetcher)}
}

func (p *TokenParser) ParseToken(ctx context.Context, accessToken string) (*Token, error) {
	token, err := jwt.Parse(accessToken, p.keyFunc(ctx))
	if err != nil {
		return nil, err
	}

	return NewToken(token), nil
}

func (p *TokenParser) keyFunc(ctx context.Context) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		cert, err := p.certificateFetcher.FetchPEMCertificate(ctx, NewToken(token))
		if err != nil {
			return nil, err
		}

		return jwt.ParseRSAPublicKeyFromPEM(cert)
	}
}

func RetrieveTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingMetadata
	}

	bearer := md.Get(AuthorizationKey)
	if len(bearer) <= 0 {
		return "", ErrMissingAuthorization
	}

	return strings.TrimPrefix(bearer[0], authorizationPrefix), nil
}

func RetrieveTokenFromRequest(r *http.Request) (string, error) {
	bearer := r.Header.Get(AuthorizationKey)
	if len(bearer) <= 0 {
		return "", ErrMissingAuthorization
	}

	return strings.TrimPrefix(bearer, authorizationPrefix), nil
}

func RetrieveTokenFromPubSubMessageAttribute(r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var payload struct {
		Message struct {
			Attributes map[string]string `json:"attributes"`
		} `json:"message"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	token, ok := payload.Message.Attributes[AuthorizationKey]
	if !ok {
		return "", ErrMissingAuthorization
	}

	return token, nil
}
