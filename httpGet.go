package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// we can format this as however we would like once we figure out how we set our paths.
	// using httpbin for now.
	name := "John Doe"
	occupation := "gardener"
	// https://httpbin.org/get?name=John+Doe&occupation=gardener

	params := "name=" + url.QueryEscape(name) + "&" +
		"occupation=" + url.QueryEscape(occupation)
	path := fmt.Sprintf("https://httpbin.org/get?%s", params)

	// creating a GET reuqest to the httpbin.org with whatever path given
	resp, err := http.Get(path)
	// checking for error
	if err != nil {
		log.Fatal(err)
	}
	// closing the response body
	defer resp.Body.Close()
	// reading content
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	// print content in console
	fmt.Println(string(body))
}
