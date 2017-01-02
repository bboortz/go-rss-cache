package main

import (
	"github.com/julienschmidt/httprouter"
	"restcache"
)

func NewRouter() *httprouter.Router {
	router := restcache.NewRouter()
	router.PUT("/item", HandlerItemCreate)
	router.POST("/item", HandlerItemCreate)
	router.GET("/item/:name", HandlerItemRead)
	router.GET("/items", HandlerItemsRead)

	return router
}
