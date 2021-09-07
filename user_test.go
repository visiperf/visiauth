package visiauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerOrganizationRoles(t *testing.T) {
	c := Customer{
		roles: map[string]string{
			"2e141238-9315-4b7f-aaa0-55b9a00c6cdb": OrganizationRoleOwner.String(),
			"d1845acc-debb-4f71-bda4-eeb8f9a91495": OrganizationRoleStandard.String(),
		},
	}

	tests := []struct {
		name           string
		customer       Customer
		organizationId string
		roles          []string
	}{{
		name:           "unknown organization",
		customer:       c,
		organizationId: "azerty",
		roles:          []string{},
	}, {
		name:           "highest role",
		customer:       c,
		organizationId: "2e141238-9315-4b7f-aaa0-55b9a00c6cdb",
		roles: []string{
			OrganizationRoleOwner.String(),
			OrganizationRoleManager.String(),
			OrganizationRoleBuyer.String(),
			OrganizationRoleStandard.String(),
		},
	}, {
		name:           "lowest role",
		customer:       c,
		organizationId: "d1845acc-debb-4f71-bda4-eeb8f9a91495",
		roles: []string{
			OrganizationRoleStandard.String(),
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.roles, test.customer.OrganizationRoles(test.organizationId))
		})
	}
}

func TestCustomerHasOneOfOrganizationRoles(t *testing.T) {
	c := Customer{
		roles: map[string]string{
			"2e141238-9315-4b7f-aaa0-55b9a00c6cdb": OrganizationRoleOwner.String(),
			"d1845acc-debb-4f71-bda4-eeb8f9a91495": OrganizationRoleStandard.String(),
		},
	}

	tests := []struct {
		name           string
		customer       Customer
		organizationId string
		roles          []string
		res            bool
	}{{
		name:           "unknown organization",
		customer:       c,
		organizationId: "azerty",
		roles: []string{
			OrganizationRoleBuyer.String(),
		},
		res: false,
	}, {
		name:           "no roles specified",
		customer:       c,
		organizationId: "2e141238-9315-4b7f-aaa0-55b9a00c6cdb",
		roles:          nil,
		res:            false,
	}, {
		name:           "role is not present",
		customer:       c,
		organizationId: "d1845acc-debb-4f71-bda4-eeb8f9a91495",
		roles: []string{
			OrganizationRoleBuyer.String(),
		},
		res: false,
	}, {
		name:           "role is present",
		customer:       c,
		organizationId: "2e141238-9315-4b7f-aaa0-55b9a00c6cdb",
		roles: []string{
			OrganizationRoleBuyer.String(),
		},
		res: true,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.res, test.customer.HasOneOfOrganizationRoles(test.organizationId, test.roles...))
		})
	}
}

func TestEmployeeHasOneOfRoles(t *testing.T) {
	e := Employee{
		roles: []string{
			RoleExpert.String(),
		},
	}

	tests := []struct {
		name     string
		employee Employee
		roles    []string
		res      bool
	}{{
		name: "user without roles",
		employee: Employee{
			roles: nil,
		},
		roles: []string{
			RoleExpert.String(),
		},
		res: false,
	}, {
		name:     "no roles specified",
		employee: e,
		roles:    nil,
		res:      false,
	}, {
		name:     "role is not present",
		employee: e,
		roles: []string{
			RoleAdmin.String(),
		},
		res: false,
	}, {
		name:     "role is present",
		employee: e,
		roles: []string{
			RoleExpert.String(),
		},
		res: true,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.res, test.employee.HasOneOfRoles(test.roles...))
		})
	}
}
