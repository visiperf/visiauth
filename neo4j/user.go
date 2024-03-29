package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/visiperf/visiauth/v3"
	"github.com/visiperf/visiauth/v3/errors"
)

type userInfo struct {
	LegacyID string
	Name     string
	Email    string
}

type UserRepository struct {
	config neo4j.SessionConfig
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		config: neo4j.SessionConfig{},
	}
}

func (r *UserRepository) FetchUserByID(ctx context.Context, userID string, scopes []string) (*visiauth.User, error) {
	driver, err := neo4j.NewDriver(env.Neo4j.Uri, neo4j.BasicAuth(env.Neo4j.User, env.Neo4j.Password, ""), func(c *neo4j.Config) {
		c.Log = neo4j.ConsoleLogger(neo4j.ERROR)
	})
	if err != nil {
		return nil, errors.Internal(err)
	}
	defer driver.Close()

	session := driver.NewSession(r.config)
	defer session.Close()

	info, err := r.fetchUserInfo(ctx, session, userID)
	if err != nil {
		return nil, err
	}

	organizationsRole, organizationLegacyIDs, err := r.fetchUserOrganizations(ctx, session, userID)
	if err != nil {
		return nil, err
	}

	return visiauth.NewUser(userID, info.Name, info.Email, info.LegacyID, scopes, organizationsRole, organizationLegacyIDs), nil
}

func (r *UserRepository) fetchUserInfo(_ context.Context, session neo4j.Session, userID string) (userInfo, error) {
	res, err := session.Run(`
		match (u:User {user_id: $user_id}) return u.legacy_id as legacy_id, u.name as name, u.email as email
	`, map[string]interface{}{
		"user_id": userID,
	})
	if err != nil {
		return userInfo{}, errors.Internal(err)
	}

	rec, err := res.Single()
	if err != nil {
		return userInfo{}, errors.Internal(err)
	}

	return userInfo{
		LegacyID: rec.Values[0].(string),
		Name:     rec.Values[1].(string),
		Email:    rec.Values[2].(string),
	}, nil
}

func (r *UserRepository) fetchUserOrganizations(_ context.Context, session neo4j.Session, userID string) (map[string]string, []string, error) {
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
		return nil, nil, errors.Internal(err)
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
		return nil, nil, errors.Internal(err)
	}

	return m, s, nil
}
