package visiauth

import "github.com/bitrise-io/go-utils/sliceutil"

const (
	RoleOwner    string = "owner"
	RoleManager  string = "manager"
	RoleBuyer    string = "buyer"
	RoleStandard string = "standard"
)

// ! keep roles order
var roles = []string{
	RoleOwner,
	RoleManager,
	RoleBuyer,
	RoleStandard,
}

func indexOfRole(role string) int {
	return sliceutil.IndexOfStringInSlice(role, roles)
}
