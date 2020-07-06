package visiauth

// AuthService represent service that will be used to decode access token
type AuthService interface {
	DecodeAccessToken(token string) (*Jwt, error)
}
