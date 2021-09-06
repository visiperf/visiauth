package visiauth

const (
	OrganizationRoleOwner    OrganizationRole = "OWNER"
	OrganizationRoleManager  OrganizationRole = "MANAGER"
	OrganizationRoleBuyer    OrganizationRole = "BUYER"
	OrganizationRoleStandard OrganizationRole = "STANDARD"
)

var organizationRolesOrdered = []string{
	OrganizationRoleOwner.String(),
	OrganizationRoleManager.String(),
	OrganizationRoleBuyer.String(),
	OrganizationRoleStandard.String(),
}

type OrganizationRole string

func (or OrganizationRole) String() string {
	return string(or)
}
