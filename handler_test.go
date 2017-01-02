package main

import (
	"encoding/json"
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rsslib"
	"strings"
	"testing"
)

type WrongRssItem struct {
	Id2          int    `json:"Id"`
	Uuid2        string `json:"Uuid"`
	Channel2     string `json:"Channel"`
	Title2       string `json:"Title"`
	Link2        string `json:"Link"`
	Description2 string `json:"Description"`
	Thumbnail2   string `json:"Thumbnail"`
	PublishDate2 string `json:"PublishDate"`
	UpdateDate2  string `json:"UpdateDate"`
}

func init() {
	addItem(rsslib.RssItem{Id: 1, Uuid: "68e9e42d-a0ba-5a4c-591b-000000000001", Channel: "TestChannel", Title: "testtitle1", Link: "http://localhost"})
	addItem(rsslib.RssItem{Id: 1, Uuid: "68e9e42d-a0ba-5a4c-591b-000000000002", Channel: "TestChannel", Title: "testtitle2", Link: "http://localhost"})
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
	body := genericRouterApiTest(t, "GET", "/item/68e9e42d-a0ba-5a4c-591b-000000000001", 200)

	bodyResponse := rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Id)
	assert.NotEmpty(bodyResponse.Uuid)
	assert.NotEmpty(bodyResponse.Channel)
	assert.NotEmpty(bodyResponse.Title)
}

func TestRouterItemReadWrongItem(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/item/68e9e42d-a0ba-5a4c-591b-unknownuud99", 404)

	bodyResponse := rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Id)
	assert.Empty(bodyResponse.Uuid)
	assert.Empty(bodyResponse.Channel)
	assert.Empty(bodyResponse.Title)

}

func TestRouterItemsRead(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/items", 200)

	bodyResponse := rsslib.RssItems{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Equal(2, len(bodyResponse))
	/*
		assert.Empty(bodyResponse.Id)
		assert.Empty(bodyResponse.Name)
		assert.False(bodyResponse.Completed)
		assert.Empty(bodyResponse.Due)
	*/
}

func TestRouterItemsCountRead(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/itemscount", 200)

	bodyResponse := RssItemCount{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Count)
	assert.Equal(int64(2), bodyResponse.Count)
	/*
		assert.Empty(bodyResponse.Id)
		assert.Empty(bodyResponse.Name)
		assert.False(bodyResponse.Completed)
		assert.Empty(bodyResponse.Due)
	*/
}

func TestRouterItemCreate(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000004", Channel: "TestChannel2", Title: "testtitle3", Link: "http://localhost"}
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

func TestRouterItemCreateWithoutUuid(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Uuid: "", Channel: "TestChannel2", Title: "testtitle3", Link: "http://localhost"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 422, strings.NewReader(requestBody))

	bodyResponse := RssItemCreated{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
	assert.NotEmpty(bodyResponse.Desc)
}

func TestRouterItemCreateWithoutChannel(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000004", Channel: "", Title: "testtitle3", Link: "http://localhost"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 422, strings.NewReader(requestBody))

	bodyResponse := RssItemCreated{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
	assert.NotEmpty(bodyResponse.Desc)
}

func TestRouterItemCreateWithoutTitle(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000004", Channel: "TestChannel2", Title: "", Link: "http://localhost"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 422, strings.NewReader(requestBody))

	bodyResponse := RssItemCreated{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
	assert.NotEmpty(bodyResponse.Desc)
}

func TestRouterItemCreateWithoutLink(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000004", Channel: "TestChannel2", Title: "testtitle3", Link: ""}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 422, strings.NewReader(requestBody))

	bodyResponse := RssItemCreated{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
	assert.NotEmpty(bodyResponse.Desc)
}

func TestRouterItemCreateNotJson(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 422, strings.NewReader("id: test"))

	bodyResponse := RssItemCreated{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
	assert.NotEmpty(bodyResponse.Desc)
}

func TestRouterItemCreateMethodNotAllowed(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Id: 1, Uuid: "68e9e42d-a0ba-5a4c-591b-000000000004", Channel: "TestChannel2", Title: "testtitle4", Link: "http://localhost"}
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
