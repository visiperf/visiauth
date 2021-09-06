package visiauth

import (
	"github.com/bitrise-io/go-utils/sliceutil"
)

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

type customer struct {
	user
	roles map[string]string
}

func (c customer) OrganizationIds() []string {
	ss := make([]string, 0, len(c.roles))
	for r := range c.roles {
		ss = append(ss, r)
	}

	return ss
}

func (c customer) OrganizationRoles(organizationId string) []string {
	idx := sliceutil.IndexOfStringInSlice(c.highestRoleInOrganization(organizationId), organizationRolesOrdered)
	if idx >= 0 {
		return organizationRolesOrdered[idx:]
	}

	return []string{}
}

func (c customer) HasOneOfOrganizationRoles(organizationId string, roles ...string) bool {
	rs := c.OrganizationRoles(organizationId)
	for _, r := range roles {
		if sliceutil.IsStringInSlice(r, rs) {
			return true
		}
	}

	return false
}

func (c customer) Roles() []string {
	return []string{}
}

func (c customer) HasOneOfRoles(roles ...string) bool {
	return false
}

func (c customer) highestRoleInOrganization(organizationId string) string {
	return c.roles[organizationId]
}
