package main

import (
	"log"
	"net/http"
)

/*
param: (http.ResponseWriter) sending responses to any connected http client.
param: (http.Request) how data is retrieved from the web request
*/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter has a Write method accept []byte(str)
	w.Write([]byte("<h1>Welcome to my web server!</h1>"))
}

/**
Any function with the signature func(http.ResponseWriter, *http.Request)
can be passed to any other function that expects the http.HandlerFunc type.
**/
func main() {
	// registering handler on the servemux.
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	// fatal() returns and logs an error if something goes wrong
	log.Fatal(http.ListenAndServe(":8080", mux))
	// servemux is a http request multiplexer. request router
}
