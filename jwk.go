package visiauth

import (
	"context"
	"encoding/json"
)

type JwkFetcher interface {
	FetchJwk(ctx context.Context, kid string) (*Jwk, error)
}

type Jwk struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func (jwk Jwk) Is(kid string) bool {
	return jwk.Kid == kid
}

func (jwk Jwk) MarshalBinary() ([]byte, error) {
	return json.Marshal(jwk)
}

func (jwk *Jwk) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, jwk)
}
