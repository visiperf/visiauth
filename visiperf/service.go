package visiperf

import (
	"fmt"

	"github.com/visiperf/visiauth"
)

type authService struct {
	secret string
}

// NewAuthService must be used to instanciate new Visiperf authentication service
func NewAuthService(secret string) visiauth.AuthService {
	return &authService{secret: secret}
}

func (as *authService) DecodeAccessToken(token string) (*visiauth.Jwt, error) {
	jwt, err := newJwtFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("token to jwt conversion error: %w", err)
	}

	if err := jwt.isValid(as.secret); err != nil {
		return nil, fmt.Errorf("jwt validation error: %w", err)
	}

	if err := jwt.isExpired(); err != nil {
		return nil, fmt.Errorf("jwt expiration error: %w", err)
	}

	return jwt.toVisiauthJwt(), nil
}
