package visiauth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJwtToken(t *testing.T) {
	alg := "RS256"
	typ := "JWT"
	userType := "customer"
	organizationRoles := map[string]string{
		"e6e4a8d5-a05d-49af-afa9-afeaeb6185d1": "OWNER",
		"a46fd11a-192b-41a1-9777-d626b857161f": "MANAGER",
		"ac79fa53-2373-4cb3-a5b1-d7a046422415": "BUYER",
		"a876aa5f-70f6-4ded-92ac-052e4e2bc211": "STANDARD",
	}
	roles := []string{
		"EXPERT",
	}
	iss := "https://dev-visiperf.eu.auth0.com/"
	sub := "auth0|611e545d1e240d006a2ca76b"
	aud := []string{
		"https://visicore.secure-vp-dev.com",
		"https://dev-visiperf.eu.auth0.com/userinfo",
	}
	iatTimestamp := 1631021164
	iat := time.Date(2021, 9, 7, 15, 26, 04, 0, time.Local)
	expTimestamp := 1631028364
	exp := time.Date(2021, 9, 7, 17, 26, 04, 0, time.Local)
	azp := "aZkESddQd5oc5heshOWuNtEBU6FvY3qL"
	scope := "openid profile email visicore:access"
	scopes := []string{
		"openid",
		"profile",
		"email",
		"visicore:access",
	}

	token := NewJwtToken(*jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"https://dev.visiperf.io/claims/user_type":          userType,
		"https://dev.visiperf.io/claims/organization_roles": organizationRoles,
		"https://dev.visiperf.io/claims/roles":              roles,
		"iss":                                               iss,
		"sub":                                               sub,
		"aud":                                               aud,
		"iat":                                               iatTimestamp,
		"exp":                                               expTimestamp,
		"azp":                                               azp,
		"scope":                                             scope,
	}), "dev.visiperf.io")

	header := token.Header()
	claims := token.Claims()

	assert.Equal(t, alg, header.Alg())
	assert.Equal(t, typ, header.Typ())

	assert.Equal(t, iss, claims.Iss())
	assert.Equal(t, sub, claims.Sub())
	assert.Equal(t, aud, claims.Aud())
	assert.Equal(t, iat, claims.Iat())
	assert.Equal(t, exp, claims.Exp())
	assert.Equal(t, scope, claims.Scope())
	assert.Equal(t, scopes, claims.Scopes())
	assert.Equal(t, organizationRoles, claims.OrganizationRoles())
	assert.Equal(t, roles, claims.Roles())
	assert.Equal(t, userType, claims.UserType())
}

func TestTypeTokenToUserConverter(t *testing.T) {
	token := func(userType string) Token {
		return NewJwtToken(*jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"https://dev.visiperf.io/claims/user_type": userType,
			"https://dev.visiperf.io/claims/organization_roles": map[string]string{
				"e6e4a8d5-a05d-49af-afa9-afeaeb6185d1": "OWNER",
				"a46fd11a-192b-41a1-9777-d626b857161f": "MANAGER",
				"ac79fa53-2373-4cb3-a5b1-d7a046422415": "BUYER",
				"a876aa5f-70f6-4ded-92ac-052e4e2bc211": "STANDARD",
			},
			"https://dev.visiperf.io/claims/roles": []string{
				"EXPERT",
			},
			"iss": "https://dev-visiperf.eu.auth0.com/",
			"sub": "auth0|611e545d1e240d006a2ca76b",
			"aud": []string{
				"https://visicore.secure-vp-dev.com",
				"https://dev-visiperf.eu.auth0.com/userinfo",
			},
			"iat":   1631021164,
			"exp":   1631028364,
			"azp":   "aZkESddQd5oc5heshOWuNtEBU6FvY3qL",
			"scope": "openid profile email visicore:access",
		}), "dev.visiperf.io")
	}

	tests := []struct {
		name      string
		converter *TypeTokenToUserConverter
		token     Token
		user      User
	}{{
		name:      "invalid user type",
		converter: NewTypeTokenToUserConverter(),
		token:     token("unknown"),
		user:      nil,
	}, {
		name:      "type customer",
		converter: NewTypeTokenToUserConverter(),
		token:     token("customer"),
		user:      Customer{},
	}, {
		name:      "type employee",
		converter: NewTypeTokenToUserConverter(),
		token:     token("employee"),
		user:      Employee{},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.IsType(t, test.user, test.converter.ConvertTokenToUser(test.token))
		})
	}
}
