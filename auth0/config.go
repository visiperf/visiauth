package auth0

import (
	"log"

	"github.com/vrischmann/envconfig"
)

var env struct {
	Auth0 struct {
		Domain string
	}
}

func initConfig() {
	if err := envconfig.Init(&env); err != nil {
		log.Fatal(err)
	}
}
