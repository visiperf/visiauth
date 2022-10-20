package visiauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/visiperf/visiauth/v3/errors"
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

func (s *Service) DecodeAccessToken(ctx context.Context, accessToken string) (Authenticable, error) {
	token, err := s.tokenParser.ParseToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	switch t := token.(type) {
	case *UserToken:
		return s.user(ctx, t)
	case *MachineToken:
		return s.app(ctx, t)
	}

	return nil, errors.Internal(fmt.Errorf("unknown token type"))
}

func (s *Service) app(ctx context.Context, token *MachineToken) (Authenticable, error) {
	return NewApp(token.AppID()), nil
}

func (s *Service) user(ctx context.Context, token *UserToken) (Authenticable, error) {
	return s.userRepository.FetchUserByID(ctx, strings.Split(token.UserID(), "|")[1], token.Scopes())
}
