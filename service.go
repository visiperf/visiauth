package visiauth

import "net/http"

const defaultAuth0Domain = "dev-visiperf.eu.auth0.com"

type Service interface {
	Validate(accessToken string) error
	User(accessToken string) User
}

type Auth0Service struct {
	tokenParser TokenParser
}

func NewAuth0Service(options ...Auth0ServiceOption) *Auth0Service {
	opts := newDefaultAuth0ServiceOptions()
	for _, opt := range options {
		opt(opts)
	}

	return &Auth0Service{
		tokenParser: NewJwtTokenParser(
			NewAuth0PEMCertificateFetcher(opts.domain, &http.Client{}),
		),
	}
}

func (s *Auth0Service) Validate(accessToken string) error {
	_, err := s.tokenParser.ParseToken(accessToken)
	return err
}

func (s *Auth0Service) User(accessToken string) User {
	panic("not implemented") // TODO: Implement
}

type Auth0ServiceOptions struct {
	domain string
}

func newDefaultAuth0ServiceOptions() *Auth0ServiceOptions {
	return &Auth0ServiceOptions{
		domain: defaultAuth0Domain,
	}
}

type Auth0ServiceOption func(*Auth0ServiceOptions)

func WithDomain(domain string) Auth0ServiceOption {
	return func(opts *Auth0ServiceOptions) {
		opts.domain = domain
	}
}
