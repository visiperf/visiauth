package api

import (
	"fmt"
	"net/http"
	"time"
)

func DecodeAccessToken(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC3339)
	fmt.Fprintf(w, currentTime)
}
