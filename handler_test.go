package main

import (
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"rsslib"
	"strings"
	"testing"
	"time"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	addItem(rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000001", Channel: "TestChannel", Title: "testtitle1", Link: "http://localhost"})
	addItem(rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000002", Channel: "TestChannel", Title: "testtitle2", Link: "http://localhost"})
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
}

func TestRouterItemsCountRead(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/itemscount", 200)

	bodyResponse := ItemCount{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Count)
	assert.Equal(uint64(2), bodyResponse.Count)
}

func TestRouterItemCreate(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000004", Channel: "TestChannel2", Title: "testtitle3", Link: "http://localhost"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 201, strings.NewReader(requestBody))

	bodyResponse := ItemCUDResult{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
	assert.Equal("created", bodyResponse.Status)
}

func TestRouterItemCreateRandom(t *testing.T) {
	assert := assert.New(t)
	var uuid string = randomString(8) + "-" + randomString(4) + "-" + randomString(4) + "-" + randomString(4) + "-" + randomString(12)
	var channel string = "testchannel" + randomString(10)
	var title string = "testtitle" + randomString(10)
	var link string = "http://" + randomString(10)
	requestStruct := rsslib.RssItem{Uuid: uuid, Channel: channel, Title: title, Link: link}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 201, strings.NewReader(requestBody))

	bodyResponse := ItemCUDResult{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.Item)
	assert.NotEmpty(bodyResponse.Status)
	assert.Equal("created", bodyResponse.Status)
}

func TestRouterItemUpdate1(t *testing.T) {
	assert := assert.New(t)
	// count items initially
	body := genericRouterApiTest(t, "GET", "/itemscount", 200)

	bodyCountResponse := ItemCount{}
	if err := json.Unmarshal(body, &bodyCountResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	var initialCount uint64 = bodyCountResponse.Count

	// check id initially
	body = genericRouterApiTest(t, "GET", "/item/68e9e42d-a0ba-5a4c-591b-000000000001", 200)

	bodyReadResponse := rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyReadResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	var initialId uint64 = bodyReadResponse.Id

	// create duplicate
	var updateDate string = time.Now().UTC().String()
	requestStruct := rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000001", Channel: "TestChannel", Title: "testtitle1", Link: "http://localhost", UpdateDate: updateDate}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body = genericRouterApiTestWithRequestBody(t, "POST", "/item", 201, strings.NewReader(requestBody))

	bodyCreateResponse := ItemCUDResult{}
	if err := json.Unmarshal(body, &bodyCreateResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyCreateResponse)
	assert.NotEmpty(bodyCreateResponse.Item)
	assert.NotEmpty(bodyCreateResponse.Status)
	assert.Equal("updated", bodyCreateResponse.Status)

	// count items finally
	body = genericRouterApiTest(t, "GET", "/itemscount", 200)

	bodyCountResponse = ItemCount{}
	if err := json.Unmarshal(body, &bodyCountResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Equal(initialCount, bodyCountResponse.Count)

	// check id finally
	body = genericRouterApiTest(t, "GET", "/item/68e9e42d-a0ba-5a4c-591b-000000000001", 200)

	bodyReadResponse = rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyReadResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyReadResponse)
	assert.Equal(uint64(0), bodyReadResponse.Id)
	assert.Equal(initialId, bodyReadResponse.Id)
	assert.NotEmpty(bodyReadResponse.Uuid)
	assert.Equal("68e9e42d-a0ba-5a4c-591b-000000000001", bodyReadResponse.Uuid)
	assert.NotEmpty(bodyReadResponse.PublishDate)
	assert.NotEmpty(bodyReadResponse.UpdateDate)

	// check items finally
	body = genericRouterApiTest(t, "GET", "/items", 200)

	bodyResponse := rsslib.RssItems{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	var resp1 rsslib.RssItem = bodyResponse[0]
	assert.Equal(uint64(0), resp1.Id)
	assert.Equal("68e9e42d-a0ba-5a4c-591b-000000000001", resp1.Uuid)
	var resp2 rsslib.RssItem = bodyResponse[1]
	assert.Equal(uint64(1), resp2.Id)
	assert.Equal("68e9e42d-a0ba-5a4c-591b-000000000002", resp2.Uuid)
}

func TestRouterItemUpdate2WithoutUdateDate(t *testing.T) {
	assert := assert.New(t)
	// count items initially
	body := genericRouterApiTest(t, "GET", "/itemscount", 200)

	bodyCountResponse := ItemCount{}
	if err := json.Unmarshal(body, &bodyCountResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	var initialCount uint64 = bodyCountResponse.Count

	// check id initially
	body = genericRouterApiTest(t, "GET", "/item/68e9e42d-a0ba-5a4c-591b-000000000002", 200)

	bodyReadResponse := rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyReadResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	var initialId uint64 = bodyReadResponse.Id

	// create duplicate
	requestStruct := rsslib.RssItem{Uuid: "68e9e42d-a0ba-5a4c-591b-000000000002", Channel: "TestChannel", Title: "testtitle2", Link: "http://localhost"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body = genericRouterApiTestWithRequestBody(t, "POST", "/item", 201, strings.NewReader(requestBody))

	bodyCreateResponse := ItemCUDResult{}
	if err := json.Unmarshal(body, &bodyCreateResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyCreateResponse)
	assert.NotEmpty(bodyCreateResponse.Item)
	assert.NotEmpty(bodyCreateResponse.Status)
	assert.Equal("updated", bodyCreateResponse.Status)

	// count items finally
	body = genericRouterApiTest(t, "GET", "/itemscount", 200)

	bodyCountResponse = ItemCount{}
	if err := json.Unmarshal(body, &bodyCountResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Equal(initialCount, bodyCountResponse.Count)

	// check id finally
	body = genericRouterApiTest(t, "GET", "/item/68e9e42d-a0ba-5a4c-591b-000000000002", 200)

	bodyReadResponse = rsslib.RssItem{}
	if err := json.Unmarshal(body, &bodyReadResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyReadResponse)
	assert.Equal(uint64(1), bodyReadResponse.Id)
	assert.Equal(initialId, bodyReadResponse.Id)
	assert.NotEmpty(bodyReadResponse.Uuid)
	assert.Equal("68e9e42d-a0ba-5a4c-591b-000000000002", bodyReadResponse.Uuid)
	assert.NotEmpty(bodyReadResponse.PublishDate)
	assert.NotEmpty(bodyReadResponse.UpdateDate)

	// check items finally
	body = genericRouterApiTest(t, "GET", "/items", 200)

	bodyResponse := rsslib.RssItems{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	var resp1 rsslib.RssItem = bodyResponse[0]
	assert.Equal(uint64(0), resp1.Id)
	assert.Equal("68e9e42d-a0ba-5a4c-591b-000000000001", resp1.Uuid)
	var resp2 rsslib.RssItem = bodyResponse[1]
	assert.Equal(uint64(1), resp2.Id)
	assert.Equal("68e9e42d-a0ba-5a4c-591b-000000000002", resp2.Uuid)
}

func TestRouterItemCreateWithoutUuid(t *testing.T) {
	assert := assert.New(t)
	requestStruct := rsslib.RssItem{Uuid: "", Channel: "TestChannel2", Title: "testtitle3", Link: "http://localhost"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string(requestJson)
	body := genericRouterApiTestWithRequestBody(t, "POST", "/item", 422, strings.NewReader(requestBody))

	bodyResponse := ItemCUDResult{}
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

	bodyResponse := ItemCUDResult{}
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

	bodyResponse := ItemCUDResult{}
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

	bodyResponse := ItemCUDResult{}
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

	bodyResponse := ItemCUDResult{}
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

	bodyResponse := ItemCUDResult{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.Empty(bodyResponse.Item)
	assert.Empty(bodyResponse.Status)
}
