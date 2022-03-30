package visiauth

import (
	"github.com/bitrise-io/go-utils/sliceutil"
	"golang.org/x/exp/maps"
)

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
	OrganizationRoles(organizationId string) []string
	HasOneOfOrganizationRoles(organizationId string, roles ...string) bool
}

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

func (c Customer) Type() UserType {
	return UserTypeCustomer
}

func (c Customer) IsType(userType UserType) bool {
	return c.Type().Is(userType)
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

func (c Customer) OrganizationRoles(organizationId string) []string {
	if role, ok := c.organizations[organizationId]; ok {
		return c.rolesIncludedIn(role)
	}

	return nil
}

func (c Customer) HasOneOfOrganizationRoles(organizationId string, roles ...string) bool {
	rs := c.OrganizationRoles(organizationId)
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
