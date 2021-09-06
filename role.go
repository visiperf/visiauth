package visiauth

const (
	RoleAdmin  Role = "ADMIN"
	RoleExpert Role = "EXPERT"
)

type Role string

func (r Role) String() string {
	return string(r)
}
