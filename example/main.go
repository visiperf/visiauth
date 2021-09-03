package main

import (
	"log"
	"os"

	"github.com/visiperf/visiauth"
)

func main() {
	log.Println(visiauth.NewAuth0Service().Validate(os.Getenv("VISIAUTH_AUTH0_ACCESS_TOKEN")))
}
