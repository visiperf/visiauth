package neo4j

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var driver neo4j.Driver

func initDriver() {
	d, err := neo4j.NewDriver(env.Neo4j.Uri, neo4j.BasicAuth(env.Neo4j.User, env.Neo4j.Password, ""), func(c *neo4j.Config) {
		c.Log = neo4j.ConsoleLogger(neo4j.DEBUG)
	})
	if err != nil {
		log.Fatal(err)
	}
	driver = d
}
