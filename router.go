package main

import (
	"github.com/julienschmidt/httprouter"
	"restcache"
)

func NewRouter() *httprouter.Router {
	router := restcache.NewRouter()
	router.PUT("/service", HandlerServiceCreate)
	router.POST("/service", HandlerServiceCreate)
	router.GET("/service/:name", HandlerServiceRead)
	router.GET("/services", HandlerServicesRead)

	return router
}
