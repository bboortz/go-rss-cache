package main

import (
	"net/http"
)

func main() {
	var ipport string = ":9090"
	router := NewRouter()

	log.Info("listening on: ", ipport)
	log.Fatal(http.ListenAndServe(ipport, router))
}
