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
		match (u:User {user_id: $user_id})-[ruo:WORKS_AT|BUY_FOR|MANAGE|OWN|DEALS_WITH]->(o:Organization)
		return o.organization_id as organization_id, type(ruo) as role
		union all
		match (u:User {user_id: $user_id})-[ruh:WORKS_AT|BUY_FOR|MANAGE|OWN]->(h:Organization)-[rhn:HEAD_OF]->(n:Network)<-[ron:IN]-(o:Organization)
		return o.organization_id as organization_id, 'WORKS_AT' as role
		union all
		match (u:User {user_id: $user_id})-[ruh:DEALS_WITH]->(h:Organization)-[rhn:HEAD_OF]->(n:Network)<-[ron:IN]-(o:Organization)
		return o.organization_id as organization_id, 'DEALS_WITH' as role
	`, map[string]interface{}{
		"user_id": uid,
	})
	if err != nil {
		return nil, err
	}

	u := User{
		Id:            uid,
		Organizations: make(map[string]string),
	}
	for res.Next() {
		values := res.Record().Values
		u.Organizations[values[0].(string)] = mRoles[values[1].(string)]
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return &u, nil
}

type User struct {
	Id            string
	Organizations map[string]string
}
