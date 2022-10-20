package visiauth

const (
	RoleOwner    string = "owner"
	RoleManager  string = "manager"
	RoleBuyer    string = "buyer"
	RoleStandard string = "standard"
	RoleSupport  string = "support"
)

var mRolesIncludedInRole = map[string][]string{
	RoleOwner: {
		RoleOwner,
		RoleManager,
		RoleBuyer,
		RoleStandard,
	},
	RoleManager: {
		RoleManager,
		RoleBuyer,
		RoleStandard,
	},
	RoleBuyer: {
		RoleBuyer,
		RoleStandard,
	},
	RoleStandard: {
		RoleStandard,
	},
	RoleSupport: {
		RoleSupport,
	},
}
