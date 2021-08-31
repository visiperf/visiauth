package visiauth

type Service interface {
	Validate(accessToken string) error
	User(accessToken string) User
}
