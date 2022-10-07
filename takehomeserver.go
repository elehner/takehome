package main

import (
	"net/http"

	"github.com/elehner/takehomeserver/users"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", users.HandleUserRequest)
	http.ListenAndServe(":8080", mux)
}
