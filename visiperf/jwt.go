package visiperf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidComposition represent error when JWT is not composed correctly
const ErrInvalidComposition = "invalid jwt composition"

// ErrInvalidSecret represent error when secret is invalid
const ErrInvalidSecret = "invalid secret"

// ErrExpiredToken represent error when token is expired
const ErrExpiredToken = "token expired"

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
		return nil, fmt.Errorf("token splitting error: %w", errors.New(ErrInvalidComposition))
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
		return fmt.Errorf("jwt signature validation error: %w", errors.New(ErrInvalidSecret))
	}

	return nil
}

func (jwt *jwt) isExpired() error {
	return nil
}

func (jwt *jwt) isUnlimited() bool {
	return true
}

func (jwt *jwt) generateSignature(secret string) (string, error) {
	return "", nil
}

func (jwt *jwt) toString() (string, error) {
	return "", nil
}
