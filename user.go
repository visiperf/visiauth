package visiauth

import (
	"github.com/bitrise-io/go-utils/sliceutil"
	"golang.org/x/exp/maps"
)

type User struct {
	id            string
	scopes        []string
	organizations map[string]string
}

func NewUser(id string, scopes []string, organizations map[string]string) *User {
	return &User{id, scopes, organizations}
}

func (u User) ID() string {
	return u.id
}

func (u User) Scopes() []string {
	return u.scopes
}

func (u User) HasScope(scope string) bool {
	return sliceutil.IsStringInSlice(scope, u.scopes)
}

func (u User) OrganizationIds() []string {
	return maps.Keys(u.organizations)
}

func (u User) OrganizationRoles() map[string][]string {
	m := make(map[string][]string)
	for _, id := range u.OrganizationIds() {
		m[id] = u.RolesInOrganization(id)
	}

	return m
}

func (u User) RolesInOrganization(organizationId string) []string {
	if role, ok := u.organizations[organizationId]; ok {
		return mRolesIncludedInRole[role]
	}

	return nil
}

func (u User) HasOneOfRolesInOrganization(organizationId string, roles ...string) bool {
	rs := u.RolesInOrganization(organizationId)
	for _, role := range roles {
		if sliceutil.IsStringInSlice(role, rs) {
			return true
		}
	}

	return false
}
