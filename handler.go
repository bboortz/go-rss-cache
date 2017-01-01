package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"restcache"
	"time"
	//	"github.com/davecgh/go-spew/spew"
)

var headerContentTypeKey string = "Content-Type"
var headerContentTypeValue string = "application/json; charset=UTF-8"

/*
 * usage: curl -H "Content-Type: application/json" -d '{"name":"go-testapi"}' http://localhost:9090/service
 */
func HandlerServiceCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusCreated
	var result ServiceCreated

	var service Service
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Info("11")
		panic(err)
	}

	if err := json.Unmarshal(body, &service); err != nil {
		statusCode = 422 // 422 - Unprocessable Entity
		result = ServiceCreated{Service: service.Name, Status: "failed", Desc: err.Error()}
	} else if service.Name == "" {
		statusCode = 422 // 422 - Unprocessable Entity
		result = ServiceCreated{Service: service.Name, Status: "failed", Desc: "service name is empty"}
	} else {
		addService(service)
		result = ServiceCreated{Service: service.Name, Status: "created", Desc: ""}
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
func HandlerServiceRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result Service

	serviceName := ps.ByName("name")
	result = findService(serviceName)

	if result.Name == "" {
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
func HandlerServicesRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result Services = services

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	restcache.LogAccess(r.Method, r.RequestURI, statusCode, start)
}
