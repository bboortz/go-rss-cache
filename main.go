package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("go-router")
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

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	logAccess("Hello", r.Method, r.RequestURI, start)
}

func Test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	w.WriteHeader(http.StatusOK)
	logAccess("Test", r.Method, r.RequestURI, start)
	fmt.Fprint(w, "{'test': 1}")
}

func main() {
	var ipport string = ":9090"
//	log := logging.MustGetLogger("go-router")
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/test", Test)
    router.GET("/hello/:name", Hello)
	router.NotFound = http.FileServer(http.Dir("public"))

	log.Info("listening on: ", ipport)
    log.Fatal(http.ListenAndServe(ipport, router))
}

