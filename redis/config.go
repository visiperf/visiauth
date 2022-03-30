package redis

import (
	"log"

	"github.com/vrischmann/envconfig"
)

var env struct {
	Redis struct {
		Addr string
	}
}

func initConfig() {
	if err := envconfig.Init(&env); err != nil {
		log.Fatal(err)
	}
}
