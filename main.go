package main

import (
	"net/http"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("go-router")

func main() {
	var ipport string = ":9090"
	router := NewRouter()

	log.Info("listening on: ", ipport)
    log.Fatal(http.ListenAndServe(ipport, router))
}

