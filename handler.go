package main

import (
	"time"
    "fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
    "github.com/julienschmidt/httprouter"
//	"github.com/davecgh/go-spew/spew"
)


var headerContentTypeKey string = "Content-Type"
var headerContentTypeValue string = "application/json; charset=UTF-8"


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	logAccess(GetFunctionName(Index), r.Method, r.RequestURI, start)
	fmt.Fprint(w, "{'api': 'go-router, 'api-version': '.0.1'}")
}

func Alive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	logAccess(GetFunctionName(Index), r.Method, r.RequestURI, start)
	fmt.Fprint(w, "{'alive': true}")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	logAccess(GetFunctionName(Index), r.Method, r.RequestURI, start)
}

func ServiceCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)

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

	w.WriteHeader(http.StatusCreated)

	addService(service)

	fmt.Fprintf(w, "{'service': %s, 'status': '%s'}", service.Name, "created")
	logAccess(GetFunctionName(Index), r.Method, r.RequestURI, start)

}


func ServiceRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	logAccess(GetFunctionName(Index), r.Method, r.RequestURI, start)
}

func ServicesRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%s!\n", services)
	fmt.Println(services)
	logAccess(GetFunctionName(Index), r.Method, r.RequestURI, start)
}
