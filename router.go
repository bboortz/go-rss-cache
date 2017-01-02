package main

import (
	"github.com/julienschmidt/httprouter"
	"restcache"
)

func NewRouter() *httprouter.Router {
	router := restcache.NewRouter()
	router.PUT("/service", HandlerItemCreate)
	router.POST("/service", HandlerItemCreate)
	router.GET("/service/:name", HandlerItemRead)
	router.GET("/services", HandlerItemsRead)

	return router
}
