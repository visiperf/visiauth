package visiauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

type Token interface {
	Header() map[string]interface{}
	Kid() string
	Alg() string
	Typ() string
	Claims() map[string]interface{}
	Iss() string
	Sub() string
	Aud() string
	Iat() string
	Exp() string
	Azp() string
}

type token struct {
	*jwt.Token
}

func newToken(t *jwt.Token) *token {
	return &token{t}
}

func (t token) Header() map[string]interface{} {
	return t.Token.Header
}

func (t token) Kid() string {
	return t.Header()["kid"].(string)
}

func (t token) Alg() string {
	return t.Header()["alg"].(string)
}

func (t token) Typ() string {
	return t.Header()["typ"].(string)
}

func (t token) Claims() map[string]interface{} {
	return t.Token.Claims.(jwt.MapClaims)
}

func (t token) Iss() string {
	return t.Claims()["iss"].(string)
}

func (t token) Sub() string {
	return t.Claims()["sub"].(string)
}

func (t token) Aud() string {
	return t.Claims()["aud"].(string)
}

func (t token) Iat() string {
	return t.Claims()["iat"].(string)
}

func (t token) Exp() string {
	return t.Claims()["exp"].(string)
}

func (t token) Azp() string {
	return t.Claims()["azp"].(string)
}

type UserToken struct {
	*token
}

func NewUserToken(token *jwt.Token) *UserToken {
	return &UserToken{newToken(token)}
}

func (t UserToken) UserID() string {
	return t.Claims()["sub"].(string)
}

func (t UserToken) Scopes() []string {
	return strings.Split(t.scope(), " ")
}

func (t UserToken) scope() string {
	return t.Claims()["scope"].(string)
}

type MachineToken struct {
	*token
}

func NewMachineToken(token *jwt.Token) *MachineToken {
	return &MachineToken{newToken(token)}
}

type TokenParser struct {
	certificateFetcher *CertificateFetcher
}

func NewTokenParser(jwkFetcher JwkFetcher) *TokenParser {
	return &TokenParser{NewCertificateFetcher(jwkFetcher)}
}

func (p *TokenParser) ParseToken(ctx context.Context, accessToken string) (*UserToken, error) {
	token, err := jwt.Parse(accessToken, p.keyFunc(ctx))
	if err != nil {
		return nil, err
	}

	return NewUserToken(token), nil
}

func (p *TokenParser) keyFunc(ctx context.Context) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		cert, err := p.certificateFetcher.FetchPEMCertificate(ctx, token.Header["kid"].(string))
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

	return strings.TrimPrefix(token, authorizationPrefix), nil
}
