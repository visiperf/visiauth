package visiauth

import (
	"github.com/bitrise-io/go-utils/sliceutil"
	"golang.org/x/exp/maps"
)

type Employee struct {
	id            string
	scopes        []string
	organizations map[string]string
}

func NewEmployee(id string, scopes []string, organizations map[string]string) User {
	return &Employee{id, scopes, organizations}
}

func (e Employee) ID() string {
	return e.id
}

func (e Employee) Type() UserType {
	return UserTypeEmployee
}

func (e Employee) IsType(userType UserType) bool {
	return e.Type().Is(userType)
}

func (e Employee) Scopes() []string {
	return e.scopes
}

func (e Employee) HasScope(scope string) bool {
	return sliceutil.IsStringInSlice(scope, e.scopes)
}

func (e Employee) OrganizationIds() []string {
	return maps.Keys(e.organizations)
}

func (e Employee) OrganizationRoles() map[string][]string {
	m := make(map[string][]string)
	for id, role := range e.organizations {
		m[id] = []string{role}
	}

	return m
}

func (e Employee) RolesInOrganization(organizationId string) []string {
	return e.OrganizationRoles()[organizationId]
}

func (e Employee) HasOneOfRolesInOrganization(organizationId string, roles ...string) bool {
	rs := e.RolesInOrganization(organizationId)
	for _, role := range roles {
		if sliceutil.IsStringInSlice(role, rs) {
			return true
		}
	}

	return false
}
