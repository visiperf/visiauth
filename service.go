package visiauth

import (
	"context"
	"strings"
)

type Service struct {
	tokenParser    *TokenParser
	userRepository UserRepository
}

func NewService(jwkFetcher JwkFetcher, userRepository UserRepository) *Service {
	return &Service{
		tokenParser:    NewTokenParser(jwkFetcher),
		userRepository: userRepository,
	}
}

func (s *Service) Validate(ctx context.Context, accessToken string) error {
	_, err := s.tokenParser.ParseToken(ctx, accessToken)
	return err
}

func (s *Service) User(ctx context.Context, accessToken string) (*User, error) {
	token, err := s.tokenParser.ParseToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return s.userRepository.FetchUserByID(ctx, strings.Split(token.UserID(), "|")[1], token.Scopes())
}
