package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type UserRepository struct {
	config neo4j.SessionConfig
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		config: neo4j.SessionConfig{},
	}
}

func (r *UserRepository) FetchUserOrganizations(ctx context.Context, userID string) (map[string]string, []string, error) {
	driver, err := neo4j.NewDriver(env.Neo4j.Uri, neo4j.BasicAuth(env.Neo4j.User, env.Neo4j.Password, ""), func(c *neo4j.Config) {
		c.Log = neo4j.ConsoleLogger(neo4j.DEBUG)
	})
	if err != nil {
		return nil, nil, err
	}
	defer driver.Close()

	session := driver.NewSession(r.config)
	defer session.Close()

	res, err := session.Run(`
		match (u:User {user_id: $user_id})-[ruh:WORKS_AT|BUY_FOR|MANAGE|OWN]->(h:Organization)-[rhn:HEAD_OF]->(n:Network)<-[ron:IN]-(o:Organization)
		return o.organization_id as organization_id, 'WORKS_AT' as role, o.legacy_id as legacy_id
			union all
		match (u:User {user_id: $user_id})-[ruh:DEALS_WITH]->(h:Organization)-[rhn:HEAD_OF]->(n:Network)<-[ron:IN]-(o:Organization)
		return o.organization_id as organization_id, 'DEALS_WITH' as role, o.legacy_id as legacy_id
			union all
		match (u:User {user_id: $user_id})-[ruo:WORKS_AT|BUY_FOR|MANAGE|OWN|DEALS_WITH]->(o:Organization)
		return o.organization_id as organization_id, type(ruo) as role, o.legacy_id as legacy_id
	`, map[string]interface{}{
		"user_id": userID,
	})
	if err != nil {
		return nil, nil, err
	}

	m := make(map[string]string)
	s := make([]string, 0)
	for res.Next() {
		values := res.Record().Values
		m[values[0].(string)] = mRelationTypeRole[values[1].(string)]

		if id := values[2].(string); len(id) > 0 {
			s = append(s, id)
		}
	}

	if err := res.Err(); err != nil {
		return nil, nil, err
	}

	return m, s, nil
}
