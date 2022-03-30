package visiauth

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/visiperf/visiauth/v3/neo4j"
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
	userId := strings.ReplaceAll(claims["sub"].(string), "auth0|", "")

	organizations, err := neo4j.FetchOrganizationsByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	// TODO: implement customer or employee depending on user type
	return NewCustomer(
		userId,
		strings.Split(claims["scope"].(string), " "),
		organizations,
	), nil
}
