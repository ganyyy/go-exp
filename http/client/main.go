package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Data struct {
	Name    string `form:"name"`
	Age     int    `form:"age"`
	Address string `form:"address"`
}

var client http.Client

func init() {
	// client.Timeout = time.Second * 3
}

func doRequest(req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var decoder = json.NewDecoder(resp.Body)

	var data Data
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", data)
}

func main() {

	var body = bytes.NewBuffer(nil)
	var encoder = json.NewEncoder(body)
	encoder.Encode(Data{
		Name:    "1231",
		Age:     100,
		Address: "123211",
	})

	var request, err = http.NewRequest(http.MethodPost, "http://localhost:8899/postform2", body)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	doRequest(request)

	var urlValue = make(url.Values)
	urlValue.Add("name", "123")
	urlValue.Add("age", "10")
	urlValue.Add("address", "1231321")

	request, err = http.NewRequest(http.MethodPost, "http://localhost:8899/postform", strings.NewReader(urlValue.Encode()))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	doRequest(request)
}
