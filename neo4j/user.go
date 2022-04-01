package neo4j

import (
	"context"
)

func FetchUserByID(ctx context.Context, userID string, scopes []string) (*User, error)

type User struct {
	Id            string            `db:"user_id"`
	Type          string            `db:"type"`
	Organizations map[string]string `db:"organizations"`
}
