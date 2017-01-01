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
)

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

func TestRouterServiceRead(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/service/go-rnd", 200)

	bodyResponse := Service{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Id)
	assert.NotEmpty(bodyResponse.Name)
	assert.False(bodyResponse.Completed)
	assert.Empty(bodyResponse.Due)
}

func TestRouterServiceReadWrongService(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/service/go-rnd2", 404)

	bodyResponse := Service{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Id)
	assert.Empty(bodyResponse.Name)
	assert.False(bodyResponse.Completed)
	assert.Empty(bodyResponse.Due)
}

func TestRouterServicesRead(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/services", 200)

	bodyResponse := Services{}
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
func TestRouterServiceCreate(t *testing.T) {
	assert := assert.New(t)
	requestStruct := ServiceCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/service", 201, strings.NewReader(requestBody))

	bodyResponse := Service{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Id)
	assert.Empty(bodyResponse.Name)
	assert.False(bodyResponse.Completed)
	assert.Empty(bodyResponse.Due)
}

func TestRouterServiceCreateMethodNotAllowed(t *testing.T) {
	assert := assert.New(t)
	requestStruct := ServiceCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/service/go-rnd2", 405, strings.NewReader(requestBody))

	bodyResponse := Service{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Id)
	assert.Empty(bodyResponse.Name)
	assert.False(bodyResponse.Completed)
	assert.Empty(bodyResponse.Due)
}
