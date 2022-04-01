package visiauth

const (
	UserTypeCustomer UserType = "customer"
	UserTypeEmployee UserType = "employee"
)

type UserType string

func (ut UserType) Equals(userType UserType) bool {
	return ut.String() == userType.String()
}

func (ut UserType) Is(userType UserType) bool {
	return ut.Equals(userType)
}

func (ut UserType) String() string {
	return string(ut)
}

type User interface {
	ID() string
	Type() UserType
	IsType(userType UserType) bool
	Scopes() []string
	HasScope(scope string) bool
	OrganizationIds() []string
	OrganizationRoles() map[string]string
	RolesInOrganization(organizationId string) []string
	HasOneOfRolesInOrganization(organizationId string, roles ...string) bool
}
