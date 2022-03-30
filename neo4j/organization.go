package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func FetchOrganizationsByUser(ctx context.Context, userId string) (map[string]string, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	res, err := session.Run(`
		match (u:User {uuid: $userId})-[ruo:WORKS_AT|BUY_FOR|MANAGE|OWN]->(o:Organization)
		return o.uuid as organization_uuid, type(ruo) as role
		union all
		match (u:User {uuid: $userId})-[ruh:WORKS_AT|BUY_FOR|MANAGE|OWN]->(h:Organization)-[rhn:HEAD_OF]->(n:Network)<-[ron:IN]-(o:Organization)
		return o.uuid as organization_uuid, 'WORKS_AT' as role
	`, map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		return nil, err
	}

	os := make(map[string]string)
	for res.Next() {
		values := res.Record().Values
		os[values[0].(string)] = mRoles[values[1].(string)]
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return os, nil
}
