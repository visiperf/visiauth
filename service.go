package visiauth

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Service struct {
	tokenParser *TokenParser
}

func NewService(jwkFetcher JwkFetcher) *Service {
	return &Service{NewTokenParser(jwkFetcher)}
}

func (s *Service) Validate(ctx context.Context, accessToken string) error {
	_, err := s.tokenParser.ParseToken(ctx, accessToken)
	return err
}

func (s *Service) User(ctx context.Context, accessToken string) (User, error) {
	token, err := s.tokenParser.ParseToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	// TODO: fetch organizations by user id in neo4j

	return NewCustomer(
		claims["sub"].(string),
		strings.Split(claims["scope"].(string), " "),
		nil,
	), nil
}
