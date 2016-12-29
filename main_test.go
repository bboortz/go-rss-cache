package main

import (
	"fmt"
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
)

type ServiceInfo struct {
	Api				string `json:"api"`
	ApiVersion		string `json:"api-version"`
}

func getIndex(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	
	var serviceInfo ServiceInfo
	if err := json.Unmarshal(body, &serviceInfo); err != nil {
		fmt.Println("ERROR: ", err)
		return "", err
	}
	tag := serviceInfo.Api
	fmt.Println("tag: ", tag)
	spew.Dump(body)
	spew.Dump(serviceInfo)
	return tag, nil
}

func TestApiIndex(t *testing.T) {
	getIndex("http://localhost:9090/")

}

func TestHello(t *testing.T) {
        expectedStr := "Hello, Testing!"
        result := "Hello, Testing!"
        if result != expectedStr {
                t.Fatalf("Expected %s, got %s", expectedStr, result)
        }
}

