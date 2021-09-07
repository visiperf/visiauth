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

type Customer struct {
	user
	roles map[string]string
}

func newCustomer(user user, roles map[string]string) Customer {
	return Customer{user, roles}
}

func (c Customer) OrganizationIds() []string {
	ss := make([]string, 0, len(c.roles))
	for r := range c.roles {
		ss = append(ss, r)
	}

	return ss
}

func (c Customer) OrganizationRoles(organizationId string) []string {
	idx := sliceutil.IndexOfStringInSlice(c.highestRoleInOrganization(organizationId), organizationRolesOrdered)
	if idx >= 0 {
		return organizationRolesOrdered[idx:]
	}

	return []string{}
}

func (c Customer) HasOneOfOrganizationRoles(organizationId string, roles ...string) bool {
	rs := c.OrganizationRoles(organizationId)
	for _, r := range roles {
		if sliceutil.IsStringInSlice(r, rs) {
			return true
		}
	}

	return false
}

func (c Customer) Roles() []string {
	return []string{}
}

func (c Customer) HasOneOfRoles(roles ...string) bool {
	return false
}

func (c Customer) highestRoleInOrganization(organizationId string) string {
	return c.roles[organizationId]
}

type Employee struct {
	user
	roles []string
}

func newEmployee(user user, roles []string) Employee {
	return Employee{user, roles}
}

func (e Employee) OrganizationIds() []string {
	return []string{}
}

func (e Employee) OrganizationRoles(organizationId string) []string {
	return []string{}
}

func (e Employee) HasOneOfOrganizationRoles(organizationId string, roles ...string) bool {
	return false
}

func (e Employee) Roles() []string {
	return e.roles
}

func (e Employee) HasOneOfRoles(roles ...string) bool {
	for _, r := range roles {
		if sliceutil.IsStringInSlice(r, e.roles) {
			return true
		}
	}

	return false
}
