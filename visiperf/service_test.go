package visiperf

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visiperf/visiauth"
)

const secret = "azerty"

func TestDecodeAccessToken(t *testing.T) {
	s := NewAuthService(secret)

	var tests = []struct {
		Message string
		Token   string
		Jwt     *visiauth.Jwt
		Err     error
	}{{
		Message: "err: empty token",
		Token:   "",
		Jwt:     nil,
		Err:     ErrInvalidComposition,
	}, {
		Message: "err: jwt not in three parts",
		Token:   "a.b",
		Jwt:     nil,
		Err:     ErrInvalidComposition,
	}, {
		Message: "err: jwt not in three parts",
		Token:   "a.b.c.d",
		Jwt:     nil,
		Err:     ErrInvalidComposition,
	}, {
		Message: "err: invalid base64 content",
		Token:   "a.b.c",
		Jwt:     nil,
		Err:     base64.CorruptInputError(0),
	}, {
		Message: "err: bad secret",
		Token:   "eyJhbGciOiJTSEE1MTIiLCJ0eXAiOiJKV1QifQ==.eyJpYXQiOiIyMDIwLTA1LTI5IDExOjIxOjIwIiwiZXhwIjoiMjAyMC0wNS0yOSAxNDoyMToyMCIsInN1YiI6eyJ1c2VySWQiOjE2ODQsImdyb3VwSWQiOjEzNDd9LCJyb2xlcyI6WyJST0xFX0NMSUVOVCJdfQ==.56673c964ad2bd24f8cc1be2d0b65d45a25882eaf172247c9a1fac75222018699a542259987e8802c46ebaa00f7866953f8877f15393913b017c4cebe34c4b59",
		Jwt:     nil,
		Err:     ErrInvalidSecret,
	}, {
		Message: "err: expired token",
		Token:   "eyJhbGciOiJTSEE1MTIiLCJ0eXAiOiJKV1QifQ==.eyJpYXQiOiIyMDIwLTA1LTI5IDExOjIxOjIwIiwiZXhwIjoiMjAxOS0wNS0yOSAxNDoyMToyMCIsInN1YiI6eyJ1c2VySWQiOjE2ODQsImdyb3VwSWQiOjEzNDd9LCJyb2xlcyI6WyJST0xFX0NMSUVOVCJdfQ==.83e7a1a08f41a994cc6bdc246cb7510cc7df61ad814abffd4e8e3bcf8db5716f2a9edad076303d46efb8ffd04782e3e563427a135b017dab13667170002e4068",
		Jwt:     nil,
		Err:     ErrExpiredToken,
	}, {
		Message: "ok: unlimited token",
		Token:   "eyJhbGciOiJTSEE1MTIiLCJ0eXAiOiJKV1QifQ==.eyJpYXQiOiIyMDIwLTA1LTI5IDExOjIxOjIwIiwiZXhwIjoiMCIsInN1YiI6eyJ1c2VySWQiOjE2ODQsImdyb3VwSWQiOjEzNDd9LCJyb2xlcyI6WyJST0xFX0NMSUVOVCJdfQ==.e4d13549c863559299f4f86c5fbdbb502f285e8618a540f83fe3dcbaec839eda58d010d463aea7e258daf4d8af39fa775dc158dd65ace5cec8541bc4339cf567",
		Jwt: &visiauth.Jwt{
			Header: visiauth.Header{
				Alg: "SHA512",
				Typ: "JWT",
			},
			Payload: visiauth.Payload{
				Iat: "2020-05-29 11:21:20",
				Exp: "0",
				Sub: visiauth.Sub{
					UserID:    1684,
					CompanyID: 1347,
				},
				Roles: []string{"ROLE_CLIENT"},
			},
			Signature: "e4d13549c863559299f4f86c5fbdbb502f285e8618a540f83fe3dcbaec839eda58d010d463aea7e258daf4d8af39fa775dc158dd65ace5cec8541bc4339cf567",
		},
		Err: nil,
	}, {
		Message: "ok",
		Token:   "eyJhbGciOiJTSEE1MTIiLCJ0eXAiOiJKV1QifQ==.eyJpYXQiOiIyMDIwLTA1LTI5IDExOjIxOjIwIiwiZXhwIjoiMjAyNS0wNS0yOSAxMToyMToyMCIsInN1YiI6eyJ1c2VySWQiOjE2ODQsImdyb3VwSWQiOjEzNDd9LCJyb2xlcyI6WyJST0xFX0NMSUVOVCJdfQ==.2257515cb6f14a2164bc1e5dd1556d2eaa0863257a0ce6af5f1f57812619a3dbd970561a5a10bd5e1dd23dee680b3667174808ec1fda86f658206d8a9c87125b",
		Jwt: &visiauth.Jwt{
			Header: visiauth.Header{
				Alg: "SHA512",
				Typ: "JWT",
			},
			Payload: visiauth.Payload{
				Iat: "2020-05-29 11:21:20",
				Exp: "2025-05-29 11:21:20",
				Sub: visiauth.Sub{
					UserID:    1684,
					CompanyID: 1347,
				},
				Roles: []string{"ROLE_CLIENT"},
			},
			Signature: "2257515cb6f14a2164bc1e5dd1556d2eaa0863257a0ce6af5f1f57812619a3dbd970561a5a10bd5e1dd23dee680b3667174808ec1fda86f658206d8a9c87125b",
		},
		Err: nil,
	}}

	for _, test := range tests {
		jwt, err := s.DecodeAccessToken(test.Token)

		// error testing
		if test.Err != nil {
			assert.True(t, errors.Is(err, test.Err), test.Message)
		} else {
			assert.Nil(t, err, test.Message)
		}

		// jwt testing
		assert.Equal(t, test.Jwt, jwt, test.Message)
	}
}
