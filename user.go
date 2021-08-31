package visiauth

type User interface {
	Id() string
	Permissions() []Permission
	HasPermission(permission Permission) bool
	OrganizationIds() []string
	OrganizationRoles(organizationId string) []OrganizationRole
	HasOneOfOrganizationRoles(organizationId string, roles ...OrganizationRole) bool
	Roles() []Role
	HasOneOfRoles(roles ...Role) bool
}
