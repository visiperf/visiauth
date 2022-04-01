package neo4j

import (
	"context"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func FetchUserByID(ctx context.Context, userId string) (*User, error) {
	uid := strings.Split(userId, "|")[1]

	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	res, err := session.Run(`
		match (u:User {user_id: $userId})-[ruo:WORKS_AT|BUY_FOR|MANAGE|OWN|DEALS_WITH]->(o:Organization)
		return o.uuid as organization_uuid, type(ruo) as role, u.user_type as type
		union all
		match (u:User {user_id: $userId})-[ruh:WORKS_AT|BUY_FOR|MANAGE|OWN|DEALS_WITH]->(h:Organization)-[rhn:HEAD_OF]->(n:Network)<-[ron:IN]-(o:Organization)
		return o.uuid as organization_uuid, 'WORKS_AT' as role, u.user_type as type
	`, map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		return nil, err
	}

	u := User{Id: uid}
	for res.Next() {
		values := res.Record().Values
		u.Organizations[values[0].(string)] = mRoles[values[1].(string)]
		u.Type = values[2].(string)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return &u, nil
}

type User struct {
	Id            string
	Type          string
	Organizations map[string]string
}
