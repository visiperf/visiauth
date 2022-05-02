package visiauth

import (
	"context"

	"github.com/bitrise-io/go-utils/sliceutil"
	"golang.org/x/exp/maps"
)

type UserRepository interface {
	FetchUserByID(ctx context.Context, userID string, scopes []string) (*User, error)
}

type User struct {
	id                    string
	legacyID              string
	scopes                []string
	organizationsRole     map[string]string
	organizationLegacyIDs []string
}

func NewUser(id, legacyID string, scopes []string, organizationsRole map[string]string, organizationLegacyIDs []string) *User {
	return &User{id, legacyID, scopes, organizationsRole, organizationLegacyIDs}
}

func (u User) ID() string {
	return u.id
}

func (u User) LegacyID() string {
	return u.legacyID
}

func (u User) Scopes() []string {
	return u.scopes
}

func (u User) HasScope(scope string) bool {
	return sliceutil.IsStringInSlice(scope, u.scopes)
}

func (u User) OrganizationIds() []string {
	return maps.Keys(u.organizationsRole)
}

func (u User) OrganizationLegacyIds() []string {
	return u.organizationLegacyIDs
}

// TODO: User.OrganizationRoles will be renamed User.OrganizationsRoles in futur major version
func (u User) OrganizationRoles() map[string][]string {
	m := make(map[string][]string)
	for _, id := range u.OrganizationIds() {
		m[id] = u.RolesInOrganization(id)
	}

	return m
}

// TODO: User.RolesInOrganization will be renamed User.OrganizationRoles in futur major version
func (u User) RolesInOrganization(organizationId string) []string {
	if role, ok := u.organizationsRole[organizationId]; ok {
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

func (u User) HighestRoleInOrganizations() map[string]string {
	return u.organizationsRole
}

func (u User) HighestRoleInOrganization(organizationId string) string {
	return u.organizationsRole[organizationId]
}
