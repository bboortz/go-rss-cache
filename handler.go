package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"restcache"
	"rsslib"
	"time"
	//	"github.com/davecgh/go-spew/spew"
)

var headerContentTypeKey string = "Content-Type"
var headerContentTypeValue string = "application/json; charset=UTF-8"

/*
 * usage: curl -H "Content-Type: application/json" -d '{"name":"go-testapi"}' http://localhost:9090/service
 */
func HandlerItemCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusCreated
	var result RssItemCreated

	var service rsslib.RssItem
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Info("11")
		panic(err)
	}

	if err := json.Unmarshal(body, &service); err != nil {
		log.Info("err 422")
		statusCode = 422 // 422 - Unprocessable Entity
		result = RssItemCreated{Item: service.Title, Status: "failed", Desc: err.Error()}
	} else if service.Title == "" {
		log.Info("err 422 b")
		statusCode = 422 // 422 - Unprocessable Entity
		statusCode = 422 // 422 - Unprocessable Entity
		result = RssItemCreated{Item: service.Title, Status: "failed", Desc: "service name is empty"}
	} else {
		addItem(service)
		result = RssItemCreated{Item: service.Title, Status: "created", Desc: ""}
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	restcache.LogAccess(r.Method, r.RequestURI, statusCode, start)
}

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/services/:name
 */
func HandlerItemRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result rsslib.RssItem

	serviceName := ps.ByName("name")
	result = findItem(serviceName)

	if result.Title == "" {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(restcache.JsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	restcache.LogAccess(r.Method, r.RequestURI, statusCode, start)
}

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/services
 */
func HandlerItemsRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result rsslib.RssItems = rssItems

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	restcache.LogAccess(r.Method, r.RequestURI, statusCode, start)
}
