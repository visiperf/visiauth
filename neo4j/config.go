package neo4j

import (
	"log"

	"github.com/vrischmann/envconfig"
)

var env struct {
	Neo4j struct {
		Uri      string
		User     string
		Password string
	}
}

func initConfig() {
	if err := envconfig.Init(&env); err != nil {
		log.Fatal(err)
	}
}
