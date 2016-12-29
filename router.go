package main

import (
	"net/http"
    "github.com/julienschmidt/httprouter"
)


func NewRouter() *httprouter.Router {
//	log := logging.MustGetLogger("go-router")
    router := httprouter.New()
    router.GET("/", IndexRead)
    router.GET("/alive", AliveRead)
    router.PUT("/service", ServiceCreate)
    router.POST("/service", ServiceCreate)
	router.GET("/service/:name", ServiceRead)
	router.GET("/services", ServicesRead)
	router.NotFound = http.FileServer(http.Dir("public"))

	return router
}

