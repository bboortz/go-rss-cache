package main

import (
	"fmt"
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/op/go-logging"
	//"github.com/davecgh/go-spew/spew"
)

func getIndex(url string) (string, error) {
	var log = logging.MustGetLogger("test-getIndex")

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	
	var api Api
	if err := json.Unmarshal(body, &api); err != nil {
		fmt.Println("ERROR: ", err)
		return "", err
	}
	log.Infof( "api: %s (%s)", api.ApiName, api.ApiVersion )
	trace()
	fmt.Println( getMethodName() )
	return api.ApiName, nil
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

