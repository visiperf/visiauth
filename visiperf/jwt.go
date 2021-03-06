package visiperf

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/visiperf/visiauth"
)

// ErrInvalidComposition represent error when JWT is not composed correctly
var ErrInvalidComposition = errors.New("invalid jwt composition")

// ErrInvalidSecret represent error when secret is invalid
var ErrInvalidSecret = errors.New("invalid secret")

// ErrExpiredToken represent error when token is expired
var ErrExpiredToken = errors.New("token expired")

type jwt struct {
	Header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	Payload struct {
		Iat string `json:"iat"`
		Exp string `json:"exp"`
		Sub struct {
			UserID    int64 `json:"userId"`
			CompanyID int64 `json:"groupId"`
		} `json:"sub"`
		Roles []string `json:"roles"`
	}
	Signature string
}

func newJwtFromToken(token string) (*jwt, error) {
	var jwt jwt

	ss := strings.Split(token, ".")
	if len(ss) != 3 {
		return nil, fmt.Errorf("token splitting error: %w", ErrInvalidComposition)
	}

	// Jwt.Header
	hDec, err := base64.StdEncoding.DecodeString(ss[0])
	if err != nil {
		return nil, fmt.Errorf("base64 decode header error: %w", err)
	}

	if err := json.Unmarshal(hDec, &jwt.Header); err != nil {
		return nil, fmt.Errorf("json unmarshal header error: %w", err)
	}

	// Jwt.Payload
	pDec, err := base64.StdEncoding.DecodeString(ss[1])
	if err != nil {
		return nil, fmt.Errorf("base64 decode payload error: %w", err)
	}

	if err := json.Unmarshal(pDec, &jwt.Payload); err != nil {
		return nil, fmt.Errorf("json unmarshal payload error: %w", err)
	}

	// Jwt.Signature
	jwt.Signature = string(ss[2])

	return &jwt, nil
}

func (jwt *jwt) isValid(secret string) error {
	s, err := jwt.generateSignature(secret)
	if err != nil {
		return fmt.Errorf("jwt signature generation error: %w", err)
	}

	if jwt.Signature != s {
		return ErrInvalidSecret
	}

	return nil
}

func (jwt *jwt) isExpired() error {
	if jwt.isUnlimited() {
		return nil
	}

	// @todo: do not use location to check jwt expiration
	loc, _ := time.LoadLocation("Europe/Paris")

	exp, err := time.ParseInLocation("2006-01-02 15:04:05", jwt.Payload.Exp, loc)
	if err != nil {
		return fmt.Errorf("jwt expiration parsing error: %w", err)
	}

	if exp.Before(time.Now()) {
		return ErrExpiredToken
	}

	return nil
}

func (jwt *jwt) isUnlimited() bool {
	return jwt.Payload.Exp == "0"
}

func (jwt *jwt) generateSignature(secret string) (string, error) {
	sp, err := jwt.toString()
	if err != nil {
		return "", fmt.Errorf("jwt to string conversion error: %w", err)
	}

	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(sp))

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func (jwt *jwt) toString() (string, error) {
	bh, err := json.Marshal(&jwt.Header)
	if err != nil {
		return "", fmt.Errorf("json marshal header error: %w", err)
	}

	bp, err := json.Marshal(&jwt.Payload)
	if err != nil {
		return "", fmt.Errorf("json marshal payload error: %w", err)
	}

	return strings.Join([]string{
		base64.StdEncoding.EncodeToString(bh),
		base64.StdEncoding.EncodeToString(bp),
	}, "."), nil
}

func (jwt *jwt) toVisiauthJwt() *visiauth.Jwt {
	return &visiauth.Jwt{
		Header: visiauth.Header{
			Alg: jwt.Header.Alg,
			Typ: jwt.Header.Typ,
		},
		Payload: visiauth.Payload{
			Iat: jwt.Payload.Iat,
			Exp: jwt.Payload.Exp,
			Sub: visiauth.Sub{
				UserID:    jwt.Payload.Sub.UserID,
				CompanyID: jwt.Payload.Sub.CompanyID,
			},
			Roles: jwt.Payload.Roles,
		},
		Signature: jwt.Signature,
	}
}
