package visiauth

type Jwks struct {
	Keys []*Jwk `json:"keys"`
}
