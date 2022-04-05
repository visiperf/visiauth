package visiauth

import (
	"context"

	"github.com/visiperf/visiauth/v2/neo4j"
)

type Service struct {
	tokenParser *TokenParser
}

func NewService(jwkFetcher JwkFetcher) *Service {
	return &Service{
		tokenParser: NewTokenParser(jwkFetcher),
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

	user, err := neo4j.FetchUserByID(ctx, token.UserID())
	if err != nil {
		return nil, err
	}

	return NewUser(user.Id, token.Scopes(), user.Organizations), nil
}
