package main

import (
	"encoding/json"
	//"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"restcache"
	"rsslib"
	"strconv"
	"time"
)

var headerContentTypeKey string = "Content-Type"
var headerContentTypeValue string = "application/json; charset=UTF-8"

/*
 * usage: curl -H "Content-Type: application/json" -d '{"name":"go-testapi"}' http://localhost:9090/item
 * usage: curl -v -H "Content-Type: application/json" -d '{"id": 1, "uuid": "11", "channel":"testchannel", "title": "testtitle", "link": "http://localhost" }' http://localhost:9090/item
 */
func HandlerItemCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusCreated
	var result ItemCreated

	var item rsslib.RssItem
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Error("err in nil")
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Error("11")
		panic(err)
	}
	log.Info("no err")

	if err := json.Unmarshal(body, &item); err != nil {
		statusCode = 422 // 422 - Unprocessable Entity
		result = ItemCreated{Item: "no uuid", Status: "failed", Desc: err.Error()}
	} else if item.Uuid == "" {
		statusCode = 422 // 422 - Unprocessable Entity
		result = ItemCreated{Item: item.Title, Status: "failed", Desc: "uuid is empty"}
	} else if item.Channel == "" {
		statusCode = 422 // 422 - Unprocessable Entity
		result = ItemCreated{Item: item.Uuid, Status: "failed", Desc: "channel is empty"}
	} else if item.Title == "" {
		statusCode = 422 // 422 - Unprocessable Entity
		result = ItemCreated{Item: item.Uuid, Status: "failed", Desc: "title is empty"}
	} else if item.Link == "" {
		statusCode = 422 // 422 - Unprocessable Entity
		result = ItemCreated{Item: item.Uuid, Status: "failed", Desc: "link is empty"}
	} else {
		addItem(item)
		result = ItemCreated{Item: item.Uuid, Status: "created", Desc: strconv.FormatUint(item.Id, 10)}
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	restcache.LogAccess(r.Method, r.RequestURI, statusCode, start)
}

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/item/:name
 */
func HandlerItemRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	var result rsslib.RssItem

	itemUuid := ps.ByName("uuid")
	result = findItem(itemUuid)

	if result.Uuid == "" {
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
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/items
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

/*
 * usage: curl -H "Content-Type: application/json" http://localhost:9090/itemscount
 */
func HandlerItemsCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start := time.Now()
	w.Header().Set(headerContentTypeKey, headerContentTypeValue)
	var statusCode int = http.StatusOK
	result := ItemCount{Count: rssItems.Len()}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
	restcache.LogAccess(r.Method, r.RequestURI, statusCode, start)
}
