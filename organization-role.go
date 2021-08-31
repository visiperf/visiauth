package visiauth

const (
	OrganizationRoleOwner    OrganizationRole = "OWNER"
	OrganizationRoleManager  OrganizationRole = "MANAGER"
	OrganizationRoleBuyer    OrganizationRole = "BUYER"
	OrganizationRoleStandard OrganizationRole = "STANDARD"
)

type OrganizationRole string
