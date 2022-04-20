package auth0

import (
	"log"

	"github.com/visiperf/visiauth/v2/config"
)

func init() {
	if err := config.Init(&env); err != nil {
		log.Fatal(err)
	}
}
