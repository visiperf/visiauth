package neo4j

import "github.com/visiperf/visiauth/v3"

var mRelationTypeRole = map[string]string{
	"WORKS_AT":   visiauth.RoleStandard,
	"BUY_FOR":    visiauth.RoleBuyer,
	"MANAGE":     visiauth.RoleManager,
	"OWN":        visiauth.RoleOwner,
	"DEALS_WITH": visiauth.RoleSupport,
}
