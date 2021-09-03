package visiauth

type Service interface {
	Validate(accessToken string) error
	User(accessToken string) User
}

type Auth0Service struct {
	domain string
}

func NewAuth0Service(domain string) *Auth0Service {
	return &Auth0Service{domain}
}

func (s *Auth0Service) Validate(accessToken string) error {
	panic("not implemented") // TODO: Implement
}

func (s *Auth0Service) User(accessToken string) User {
	panic("not implemented") // TODO: Implement
}
