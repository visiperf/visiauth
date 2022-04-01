package visiauth

type User interface {
	ID() string
	Scopes() []string
	HasScope(scope string) bool
	OrganizationIds() []string
	OrganizationRoles() map[string]string
	RolesInOrganization(organizationId string) []string
	HasOneOfRolesInOrganization(organizationId string, roles ...string) bool
}
