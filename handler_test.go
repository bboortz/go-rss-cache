package main

import (
	"fmt"
	"strings"
	"testing"
	//	"reflect"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	//	"github.com/julienschmidt/httprouter"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"rsslib"
)

type TestItemCreate struct {
	Title string `json:"Title"`
}

func init() {
	addItem(rsslib.RssItem{Channel: "TestChannel", Title: "TestTitle"})
	addItem(rsslib.RssItem{Channel: "TestChannel", Title: "TestTitle2"})
}

func genericRouterApiTest(t *testing.T, method string, url string, expectedStatusCode int) []byte {
	return genericRouterApiTestWithRequestBody(t, method, url, expectedStatusCode, nil)
}

func genericRouterApiTestWithRequestBody(t *testing.T, method string, url string, expectedStatusCode int, requestBody io.Reader) []byte {
	assert := assert.New(t)
	router := NewRouter()

	req, err := http.NewRequest(method, url, requestBody)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Nil(err)
	assert.NotNil(req)
	assert.NotNil(recorder)
	assert.Equal(expectedStatusCode, recorder.Code)

	body, err := ioutil.ReadAll(io.LimitReader(recorder.Body, 1048576))
	if err != nil {
		panic(err)
	}
	assert.NotNil(body)

	return body
}

func TestRouterItemRead(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/item/TestTitle", 200)

	bodyResponse := rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Id)
	assert.NotEmpty(bodyResponse.Channel)
	assert.NotEmpty(bodyResponse.Title)
}

func TestRouterItemReadWrongService(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/item/TestTitle100", 404)

	bodyResponse := rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Id)
	assert.Empty(bodyResponse.Channel)
	assert.Empty(bodyResponse.Title)

}

func TestRouterItemsRead(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/items", 200)

	bodyResponse := rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	/*
		assert.Empty(bodyResponse.Id)
		assert.Empty(bodyResponse.Name)
		assert.False(bodyResponse.Completed)
		assert.Empty(bodyResponse.Due)
	*/
}
func TestRouterItemCreate(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Channel: "TestChannel2", Title: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 201, strings.NewReader(requestBody))

	bodyResponse := RssItemCreated{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
}

func TestRouterItemCreateMethodNotAllowed(t *testing.T) {
	assert := assert.New(t)
	requestStruct := RssItemCreate{Channel: "TestChannel2", Title: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item/go-rnd2", 405, strings.NewReader(requestBody))

	bodyResponse := RssItemCreated{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Item)
	assert.Empty(bodyResponse.Status)
}
