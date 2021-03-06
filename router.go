package main

import (
	"github.com/bboortz/go-restcache"
	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	router := restcache.NewRouter()
	router.PUT("/item", HandlerItemCreate)
	router.POST("/item", HandlerItemCreate)
	router.GET("/item/:uuid", HandlerItemRead)
	router.GET("/items", HandlerItemsRead)
	router.GET("/itemscount", HandlerItemsCount)

	return router
}
