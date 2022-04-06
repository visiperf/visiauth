package main

import (
	"context"
	"log"

	"github.com/visiperf/visiauth/v2"
	"github.com/visiperf/visiauth/v2/neo4j"
	"github.com/visiperf/visiauth/v2/redis"
)

func main() {
	user, err := visiauth.NewService(
		redis.NewJwkFetcher(),
		neo4j.NewUserRepository(),
	).User(context.Background(), env.Visiauth.Token)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
