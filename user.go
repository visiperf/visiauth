package visiauth

import "github.com/bitrise-io/go-utils/sliceutil"

type User interface {
	Id() string
	Permissions() []string
	HasPermission(permission string) bool
	OrganizationIds() []string
	OrganizationRoles(organizationId string) []string
	HasOneOfOrganizationRoles(organizationId string, roles ...string) bool
	Roles() []string
	HasOneOfRoles(roles ...string) bool
}

type user struct {
	id          string
	permissions []string
}

func (u user) Id() string {
	return u.id
}

func (u user) Permissions() []string {
	return u.permissions
}

func (u user) HasPermission(permission string) bool {
	return sliceutil.IsStringInSlice(permission, u.permissions)
}
