package main

import (
	"log"
	"os"

	"github.com/visiperf/visiauth"
)

func main() {
	accessToken := os.Getenv("VISIAUTH_AUTH0_ACCESS_TOKEN")
	service := visiauth.NewAuth0Service()

	if err := service.Validate(accessToken); err != nil {
		log.Fatal(err)
	}

	user := service.User(accessToken)

	log.Println(user.Id())
	log.Println(user.Permissions())
	log.Println(user.Roles())

	for _, id := range user.OrganizationIds() {
		log.Println(id)
		log.Println(user.OrganizationRoles(id))
	}
}
