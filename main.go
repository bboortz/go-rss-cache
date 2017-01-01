package main

import (
	"github.com/op/go-logging"
	"net/http"
)

var log = logging.MustGetLogger("go-router")

func main() {
	var ipport string = ":9090"
	router := NewRouter()

	log.Info("listening on: ", ipport)
	log.Fatal(http.ListenAndServe(ipport, router))
}
