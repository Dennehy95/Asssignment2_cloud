package resource

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
)

func TestHandlerPostWithGetDel(t *testing.T) {
	testServer1 := httptest.NewServer(http.HandlerFunc(HandlerPost))
	testServer2 := httptest.NewServer(http.HandlerFunc(HandlerGetDel))
	defer testServer1.Close()
	defer testServer2.Close()

	form := Project{"", "testURL", "EUR", "NOK", 2.4, 5.2}

	jsonContent, err := json.Marshal(form)
	if err != nil {
		fmt.Println("Error occured during json.Marshal: ", err)
		return
	}

	response1, err := http.Post(testServer1.URL, "application/json", bytes.NewBuffer(jsonContent))
	if err != nil {
		fmt.Println("Error occured during http.Post: ", err)
		return
	}

	id, err := ioutil.ReadAll(response1.Body)
	if err != nil {
		fmt.Println("Error occured during ioutil.ReadAll: ", err)
		return
	}

	_, err = http.Get(testServer2.URL + "/" + string(id))
	if err != nil {
		fmt.Println("Error occured during http.Get: ", err)
		return
	}

	client := &http.Client{}

	response2, err := http.NewRequest(http.MethodDelete, testServer2.URL + "/" + string(id), nil)
	if err != nil {
		fmt.Println("Error occured during http.NewRequest: ", err)
		return
	}

	_, err = client.Do(response2)
	if err != nil {
		fmt.Println("Error occured during client.Do(): ", err)
	}
}

/*
func TestHandlerGet(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandlerGetDel))
	defer testServer.Close()

	id := "59f31bc86a3b2e31f4066829"


	response, err := http.Get(testServer.URL + "/" + id)
	if err != nil {
		fmt.Println("Error occured during true http.Get: ", err)
		return
	}

	fmt.Println(result)
}
*/

/*
func TestAutoTriggerCheck(t *testing.T) {
	err := AutoTriggerCheck
	//defer testServer.Close()

	//_, err = http.Get(testServer.URL)
	if err != nil {
		fmt.Println("Error occured during http.Get: ", err)
	}
}
*/

func TestFullTriggerCheck(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(FullTriggerCheck))
	defer testServer.Close()

	_, err := http.Get(testServer.URL)
	if err != nil {
		fmt.Println("Error occured during http.Get: ", err)
	}
}


func TestGetJSON(t *testing.T) {
	test := Rates{}

	err := GetJSON("https://pespiri.com/go/test.json", &test)
	if err != nil {
		fmt.Println("Error occured during http.Get: ", err)
		return
	}
}


func TestHandlerAverage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandlerAverage))
	defer testServer.Close()

	post := make(map[string]interface{})

	post["baseCurrency"]   = "EUR"
	post["targetCurrency"] = "EUR"

	jsonContent, err := json.Marshal(post)
	if err != nil {
		fmt.Print("Error occured during json.Marshal(): ", err)
		return
	}

	_, err = http.Post(testServer.URL, "application/json", bytes.NewBuffer(jsonContent))
	if err != nil {
		fmt.Print("Error occured during http.Post(): ", err)
		return
	}
}


func TestHandlerLatest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandlerLatest))
	defer testServer.Close()

	post := make(map[string]interface{})

	post["baseCurrency"]   = "EUR"
	post["targetCurrency"] = "NOK"

	jsonContent, err := json.Marshal(post)
	if err != nil {
		fmt.Print("Error occured during json.Marshal(): ", err)
		return
	}

	//response, err := http.Post(testServer.URL, "application/json", bytes.NewBuffer(jsonContent))
	_, err = http.Post(testServer.URL, "application/json", bytes.NewBuffer(jsonContent))
	if err != nil {
		fmt.Print("Error occured during http.Post(): ", err)
		return
	}

	//fmt.Print(response)
}

func TestInvoker(t *testing.T) {

}

func TestStartSession(t *testing.T) {
	_, session, err:= StartSession("mongodb://user:123@ds133465.mlab.com:33465/assignment2", "assignment2", "tottot")
	defer session.Close()

	if err != nil {
		fmt.Print("Error occured during StartSession: ", err)
		return
	}
}

/*
//WARNING! This test is unable to be finished due to the infinite loop in the function
func TestTimedTicker(t *testing.T) {
	err := TimedTicker(time.Second * 5)
	if err != nil {
		fmt.Print("Error occured during TimedTicker: ", err)
		return
	}
}
*/