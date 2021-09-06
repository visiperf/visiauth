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

type UserInstanciator interface {
	ForType() string
	InstanciateUser(token Token) User
}

type CustomerUserInstanciator struct{}

func NewCustomerUserInstanciator() *CustomerUserInstanciator {
	return &CustomerUserInstanciator{}
}

func (i *CustomerUserInstanciator) ForType() string {
	return "customer"
}

func (i *CustomerUserInstanciator) InstanciateUser(token Token) User {
	claims := token.Claims()
	return newCustomer(newUser(claims.Sub(), claims.Scopes()), claims.OrganizationRoles())
}

type EmployeeUserInstanciator struct{}

func NewEmployeeUserInstanciator() *EmployeeUserInstanciator {
	return &EmployeeUserInstanciator{}
}

func (i *EmployeeUserInstanciator) ForType() string {
	return "employee"
}

func (i *EmployeeUserInstanciator) InstanciateUser(token Token) User {
	claims := token.Claims()
	return newEmployee(newUser(claims.Sub(), claims.Scopes()), claims.Roles())
}

type user struct {
	id          string
	permissions []string
}

func newUser(id string, permissions []string) user {
	return user{id, permissions}
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

func newCustomer(user user, roles map[string]string) customer {
	return customer{user, roles}
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

type employee struct {
	user
	roles []string
}

func newEmployee(user user, roles []string) employee {
	return employee{user, roles}
}

func (e employee) OrganizationIds() []string {
	return []string{}
}

func (e employee) OrganizationRoles(organizationId string) []string {
	return []string{}
}

func (e employee) HasOneOfOrganizationRoles(organizationId string, roles ...string) bool {
	return false
}

func (e employee) Roles() []string {
	return e.roles
}

func (e employee) HasOneOfRoles(roles ...string) bool {
	for _, r := range roles {
		if sliceutil.IsStringInSlice(r, e.roles) {
			return true
		}
	}

	return false
}
