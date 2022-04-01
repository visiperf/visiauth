package visiauth

import (
	"context"

	"github.com/visiperf/visiauth/v3/neo4j"
)

type Service struct {
	tokenParser   *TokenParser
	instanciators map[UserType]func(id string, scopes []string, organizations map[string]string) User
}

func NewService(jwkFetcher JwkFetcher) *Service {
	return &Service{
		tokenParser: NewTokenParser(jwkFetcher),
		instanciators: map[UserType]func(id string, scopes []string, organizations map[string]string) User{
			UserTypeCustomer: NewCustomer,
		},
	}
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

	user, err := neo4j.FetchUserByID(ctx, token.UserID(), token.Scopes())
	if err != nil {
		return nil, err
	}

	return s.instanciators[UserType(user.Type)](user.Id, token.Scopes(), user.Organizations), nil
}
