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

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090
 */
func IndexRead(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	var api Api = Api{ApiName: "go-router", ApiVersion: "0.1",}
	if err := json.NewEncoder(w).Encode(api); err != nil {
		panic(err)
	}
	logAccess(GetFunctionName(IndexRead), r.Method, r.RequestURI, start)
}

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090
 */
func AliveRead(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	var alive Alive = Alive{Alive: true,}
	if err := json.NewEncoder(w).Encode(alive); err != nil {
		panic(err)
	}
	logAccess(GetFunctionName(IndexRead), r.Method, r.RequestURI, start)
}

/*
 * usage: curl -H "Content-Type: application/json" -d '{"name":"go-testapi"}' http://localhost:9090/service
 */
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
	logAccess(GetFunctionName(IndexRead), r.Method, r.RequestURI, start)

}


/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/services/:name
 */
func ServiceRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	serviceName := ps.ByName("name")
	s := findService(serviceName)

	if (s.Name == "") {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

	logAccess(GetFunctionName(IndexRead), r.Method, r.RequestURI, start)
}

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/services
 */
func ServicesRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(services); err != nil {
		panic(err)
	}
	logAccess(GetFunctionName(IndexRead), r.Method, r.RequestURI, start)
}
