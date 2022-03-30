package main

import (
	"log"

	"github.com/vrischmann/envconfig"
)

var env struct {
	Visiauth struct {
		Token string
	}
}

func initConfig() {
	if err := envconfig.Init(&env); err != nil {
		log.Fatal(err)
	}
}
