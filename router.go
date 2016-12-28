package main

import (
	"time"
    "fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
    "github.com/julienschmidt/httprouter"
)


var headerContentTypeKey string = "Content-Type"
var headerContentTypeValue string = "application/json; charset=UTF-8"


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	logAccess("Index", r.Method, r.RequestURI, start)
	printStack()
	fmt.Fprint(w, "{'api': 'go-router, 'api-version': '.0.1'}")
}

func Alive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	logAccess("Test", r.Method, r.RequestURI, start)
	fmt.Fprint(w, "{'alive': true}")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	logAccess("Hello", r.Method, r.RequestURI, start)
}

func ServiceCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
//	w.Header().Set(headerContentTypeKey, headerContentTypeValue)

	var service Service

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &service); err != nil {
//		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

//	w.WriteHeader(http.StatusCreated)

	addService(service)

	fmt.Fprintf(w, "{'service': %s, 'status': '%s'}", service.Name, "created")
	logAccess("ServiceCreate", r.Method, r.RequestURI, start)

	// fmt.Println(service)
	// fmt.Println(service.getJson(service))
	// fmt.Println("s1: ", service.getStruct(body))
//	fmt.Println("s: ", service.UnmarshalJSON("{'id': 1, 'name': 'lala'}"))

/*
	b, err := json.Marshal(service)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	*/

	

}


func NewRouter() *httprouter.Router {
//	log := logging.MustGetLogger("go-router")
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/alive", Alive)
    router.GET("/hello/:name", Hello)
    router.PUT("/service", ServiceCreate)
    router.POST("/service", ServiceCreate)
	router.NotFound = http.FileServer(http.Dir("public"))

	return router
}

