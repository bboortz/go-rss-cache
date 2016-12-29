package main

import (
	"time"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
    "github.com/julienschmidt/httprouter"
//	"github.com/davecgh/go-spew/spew"
)


var headerContentTypeKey string = "Content-Type"
var headerContentTypeValue string = "application/json; charset=UTF-8"

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090
 */
func IndexRead(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result Api

	w.WriteHeader(statusCode)
	result = Api{ApiName: "go-router", ApiVersion: "0.1",}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	logAccess(getMethodName(), r.Method, r.RequestURI, statusCode, start)
}

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090
 */
func AliveRead(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result Alive

	w.WriteHeader(statusCode)
	result = Alive{Alive: true,}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	logAccess(getMethodName(), r.Method, r.RequestURI, statusCode, start)
}

/*
 * usage: curl -H "Content-Type: application/json" -d '{"name":"go-testapi"}' http://localhost:9090/service
 */
func ServiceCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		panic(err)
	}
	if err := json.Unmarshal(body, &service); err != nil {
		w.Header().Set(headerContentTypeKey, headerContentTypeValue)
		w.WriteHeader(422) 
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.WriteHeader(statusCode)
	addService(service)

	result = ServiceCreated{Service: service.Name, Status: "created",}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	logAccess(getMethodName(), r.Method, r.RequestURI, statusCode, start)
}


/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/services/:name
 */
func ServiceRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result Service


	serviceName := ps.ByName("name")
	result = findService(serviceName)

	if (result.Name == "") {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}

	logAccess(getMethodName(), r.Method, r.RequestURI, statusCode, start)
}

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/services
 */
func ServicesRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result Services = services

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	logAccess(getMethodName(), r.Method, r.RequestURI, statusCode, start)
}

