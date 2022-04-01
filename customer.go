package visiauth

import (
	"github.com/bitrise-io/go-utils/sliceutil"
	"golang.org/x/exp/maps"
)

type Customer struct {
	id            string
	scopes        []string
	organizations map[string]string
}

func NewCustomer(id string, scopes []string, organizations map[string]string) *Customer {
	return &Customer{id, scopes, organizations}
}

func (c Customer) ID() string {
	return c.id
}

func (c Customer) Scopes() []string {
	return c.scopes
}

func (c Customer) HasScope(scope string) bool {
	return sliceutil.IsStringInSlice(scope, c.scopes)
}

func (c Customer) OrganizationIds() []string {
	return maps.Keys(c.organizations)
}

func (c Customer) OrganizationRoles() map[string]string {
	return c.organizations
}

func (c Customer) RolesInOrganization(organizationId string) []string {
	if role, ok := c.organizations[organizationId]; ok {
		return c.rolesIncludedIn(role)
	}

	return nil
}

func (c Customer) HasOneOfRolesInOrganization(organizationId string, roles ...string) bool {
	rs := c.RolesInOrganization(organizationId)
	for _, role := range roles {
		if sliceutil.IsStringInSlice(role, rs) {
			return true
		}
	}

	return false
}

func (c Customer) rolesIncludedIn(role string) []string {
	return roles[indexOfRole(role):]
}
