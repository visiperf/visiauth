package neo4j

import (
	"log"

	"github.com/visiperf/visiauth/v2/config"
)

func init() {
	if err := config.Init(&env); err != nil {
		log.Fatal(err)
	}
}
