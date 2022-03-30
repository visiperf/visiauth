package visiauth

import "context"

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
