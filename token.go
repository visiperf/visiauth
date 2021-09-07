package visiauth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Token interface {
	Header() Header
	Claims() Claims
	Signature() string
}

type Header interface {
	Alg() string
	Typ() string
	Kid() string
}

type Claims interface {
	Iss() string
	Sub() string
	Aud() []string
	Iat() time.Time
	Exp() time.Time
	Scope() string
	Scopes() []string
	OrganizationRoles() map[string]string
	Roles() []string
	UserType() string
}

type JwtToken struct {
	jwt.Token
	domain string
}

func NewJwtToken(token jwt.Token, domain string) JwtToken {
	return JwtToken{token, domain}
}

func (t JwtToken) Header() Header {
	return JwtHeader(t.Token.Header)
}

func (t JwtToken) Claims() Claims {
	return JwtClaims{
		domain:    t.domain,
		MapClaims: t.Token.Claims.(jwt.MapClaims),
	}
}

func (t JwtToken) Signature() string {
	return t.Token.Signature
}

type JwtHeader map[string]interface{}

func (h JwtHeader) Alg() string {
	return h["alg"].(string)
}

func (h JwtHeader) Typ() string {
	return h["typ"].(string)
}

func (h JwtHeader) Kid() string {
	return h["kid"].(string)
}

type JwtClaims struct {
	jwt.MapClaims
	domain string
}

func (c JwtClaims) Iss() string {
	return c.MapClaims["iss"].(string)
}

func (c JwtClaims) Sub() string {
	return c.MapClaims["sub"].(string)
}

func (c JwtClaims) Aud() []string {
	return c.MapClaims["aud"].([]string)
}

func (c JwtClaims) Iat() time.Time {
	return time.Unix(int64(c.MapClaims["iat"].(int)), 0)
}

func (c JwtClaims) Exp() time.Time {
	return time.Unix(int64(c.MapClaims["exp"].(int)), 0)
}

func (c JwtClaims) Scope() string {
	return c.MapClaims["scope"].(string)
}

func (c JwtClaims) Scopes() []string {
	return strings.Split(c.Scope(), " ")
}

func (c JwtClaims) OrganizationRoles() map[string]string {
	roles := make(map[string]string)
	for k, v := range c.MapClaims[c.customKey("organization_roles")].(map[string]interface{}) {
		roles[k] = v.(string)
	}

	return roles
}

func (c JwtClaims) Roles() []string {
	roles := make([]string, 0)
	for _, r := range c.MapClaims[c.customKey("roles")].([]interface{}) {
		roles = append(roles, r.(string))
	}

	return roles
}

func (c JwtClaims) UserType() string {
	return c.MapClaims[c.customKey("user_type")].(string)
}

func (c JwtClaims) customKey(key string) string {
	return fmt.Sprintf("https://%s/claims/%s", c.domain, key)
}

type TokenParser interface {
	ParseToken(accessToken string) (Token, error)
}

type JwtTokenParser struct {
	domain             string
	certificateFetcher PEMCertificateFetcher
}

func NewJwtTokenParser(domain string, certificateFetcher PEMCertificateFetcher) *JwtTokenParser {
	return &JwtTokenParser{domain, certificateFetcher}
}

func (p *JwtTokenParser) ParseToken(accessToken string) (Token, error) {
	token, err := jwt.Parse(accessToken, p.keyFunc)
	if err != nil {
		return nil, err
	}

	return NewJwtToken(*token, p.domain), nil
}

func (p *JwtTokenParser) keyFunc(token *jwt.Token) (interface{}, error) {
	cert, err := p.certificateFetcher.FetchPEMCertificate(NewJwtToken(*token, p.domain))
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(cert)
}

type TokenToUserConverter interface {
	ConvertTokenToUser(token Token) User
}

type TypeTokenToUserConverter struct {
	instanciators []UserInstanciator
}

func NewTypeTokenToUserConverter() *TypeTokenToUserConverter {
	return &TypeTokenToUserConverter{
		instanciators: []UserInstanciator{
			NewCustomerUserInstanciator(),
			NewEmployeeUserInstanciator(),
		},
	}
}

func (c *TypeTokenToUserConverter) ConvertTokenToUser(token Token) User {
	userType := token.Claims().UserType()
	for _, i := range c.instanciators {
		if i.ForType() == userType {
			return i.InstanciateUser(token)
		}
	}

	return nil
}
