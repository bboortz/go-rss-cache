package main

import (
    "github.com/julienschmidt/httprouter"
)


func NewRouter() *httprouter.Router {
//	log := logging.MustGetLogger("go-router")
    router := httprouter.New()
    router.GET("/", HandlerIndexRead)
    router.GET("/alive", HandlerAliveRead)
    router.PUT("/service", HandlerServiceCreate)
    router.POST("/service", HandlerServiceCreate)
	router.GET("/service/:name", HandlerServiceRead)
	router.GET("/services", HandlerServicesRead)
	router.NotFound = NotFoundHandler()
	router.MethodNotAllowed = MethodNotAllowedHandler()

	return router
}

