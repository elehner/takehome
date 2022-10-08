package main

import (
	"net/http"

	"github.com/elehner/takehomeserver/images"
	"github.com/elehner/takehomeserver/users"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", users.HandleUserRequest)
	mux.HandleFunc("/image", images.HandleImageRequest)
	http.ListenAndServe(":8080", mux)
}
