package main

import (
	"github.com/julienschmidt/httprouter"
	"httpcache"
)

func NewRouter() *httprouter.Router {
	router := httpcache.NewRouter()
	router.PUT("/service", HandlerServiceCreate)
	router.POST("/service", HandlerServiceCreate)
	router.GET("/service/:name", HandlerServiceRead)
	router.GET("/services", HandlerServicesRead)

	return router
}
