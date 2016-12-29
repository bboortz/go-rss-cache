package main

import (
	"fmt"
	"strings"
	"testing"
//	"reflect"
	"io"
	"io/ioutil"
	"net/http"
//	"net/http/httptest"
	"encoding/json"
//	"github.com/julienschmidt/httprouter"
//	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type mockResponseWriter struct{
	Code	int
	Body	io.Reader
}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}


type handlerStruct struct {
	handeled *bool
}

func (h handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*h.handeled = true
}

/*
func TestRouter(t *testing.T) {
	router := httprouter.New()

	routed := false
	router.Handle("GET", "/user/:name", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		routed = true
		want := httprouter.Params{httprouter.Param{"name", "gopher"}}
		if !reflect.DeepEqual(ps, want) {
			t.Fatalf("wrong wildcard values: want %v, got %v", want, ps)
		}
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user/gopher", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing failed")
	}
}
*/
func genericRouterApiMock(t *testing.T, method string, url string, expectedStatusCode int) []byte {
	return genericRouterApiTestWithRequestBody(t, method, url, expectedStatusCode, nil)
}




func genericRouterApiMockWithRequestBody(t *testing.T, method string, url string, expectedStatusCode int, requestBody io.Reader) []byte {
	assert := assert.New(t) 
	router := NewRouter()

	req, err := http.NewRequest(method, url, requestBody)
	recorder := new(mockResponseWriter)
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

func TestRouterIndexReadMock(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/", 200)

	bodyResponse := Api{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.NotEmpty(bodyResponse.ApiName)
	assert.NotEmpty(bodyResponse.ApiVersion)
}

func TestRouterAliveReadMock(t *testing.T) {
	assert := assert.New(t)
	body := genericRouterApiTest(t, "GET", "/alive", 200)

	bodyResponse := Alive{}
	if err := json.Unmarshal(body, &bodyResponse); err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.NotNil(bodyResponse)
	assert.True(bodyResponse.Alive)
}

func TestRouterServiceReadMock(t *testing.T) {
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

func TestRouterServiceReadWrongServiceMock(t *testing.T) {
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

func TestRouterServicesReadMock(t *testing.T) {
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
func TestRouterServiceCreateMock(t *testing.T) {
	assert := assert.New(t)
	requestStruct := ServiceCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string( requestJson )
	body := genericRouterApiTestWithRequestBody(t, "POST", "/service", 201, strings.NewReader(requestBody) )

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

func TestRouterServiceCreateMethodNotAllowedMock(t *testing.T) {
	assert := assert.New(t)
	requestStruct := ServiceCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string( requestJson )
	body := genericRouterApiTestWithRequestBody(t, "POST", "/service/go-rnd2", 405, strings.NewReader(requestBody) )

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


func TestRouterServiceCreateNotFoundMock(t *testing.T) {
	assert := assert.New(t)
	requestStruct := ServiceCreate{Name: "go-test"}
	requestJson, _ := json.Marshal(requestStruct)
	requestBody := string( requestJson )
	body := genericRouterApiTestWithRequestBody(t, "POST", "/notfound", 404, strings.NewReader(requestBody) )

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



