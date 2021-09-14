package visiauth

import "net/http"

const defaultAuth0Domain = "dev-visiperf.eu.auth0.com"
const defaultAuth0Namespace = "dev.visiperf.io"

type Service interface {
	Validate(accessToken string) error
	User(accessToken string) User
}

type Auth0Service struct {
	tokenParser    TokenParser
	tokenConverter TokenToUserConverter
}

/*
Create a new instance of authentication service using Auth0 provider

default domain: dev-visiperf.eu.auth0.com
*/
func NewAuth0Service(options ...Auth0ServiceOption) *Auth0Service {
	opts := newDefaultAuth0ServiceOptions()
	for _, opt := range options {
		opt(opts)
	}

	return &Auth0Service{
		tokenParser: NewJwtTokenParser(
			opts.namespace,
			NewAuth0PEMCertificateFetcher(opts.domain, &http.Client{}),
		),
		tokenConverter: NewTypeTokenToUserConverter(),
	}
}

// Parse and validate access token
func (s *Auth0Service) Validate(accessToken string) error {
	_, err := s.tokenParser.ParseToken(accessToken)
	return err
}

// Return authenticated User, nil if token is invalid
func (s *Auth0Service) User(accessToken string) User {
	token, err := s.tokenParser.ParseToken(accessToken)
	if err != nil {
		return nil
	}

	return s.tokenConverter.ConvertTokenToUser(token)
}

type Auth0ServiceOptions struct {
	domain    string
	namespace string
}

func newDefaultAuth0ServiceOptions() *Auth0ServiceOptions {
	return &Auth0ServiceOptions{
		domain:    defaultAuth0Domain,
		namespace: defaultAuth0Namespace,
	}
}

type Auth0ServiceOption func(*Auth0ServiceOptions)

/*
Use this option to change Auth0 domain
*/
func WithDomain(domain string) Auth0ServiceOption {
	return func(opts *Auth0ServiceOptions) {
		opts.domain = domain
	}
}

/*
Use this option to change Auth0 namespace
*/
func WithNamespace(namespace string) Auth0ServiceOption {
	return func(opts *Auth0ServiceOptions) {
		opts.namespace = namespace
	}
}
