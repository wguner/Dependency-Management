package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
assuming we will use a generic url scheme
[//[user:password@]host[:port]][/]path[?query][#fragment]

*/
func handler(w http.ResponseWriter, r *http.Request) {
	// name is our parameter
	keys, ok := r.URL.Query()["name"]
	name := "guest"

	if ok {

		name = keys[0]
	}
	// run curl localhost:8080/?name=Setenay in terminal
	fmt.Fprintf(w, "Hello %s!\n", name)
	// sending name as a query parameter and server responds with a message
}

func main() {

	http.HandleFunc("/", handler)
	fmt.Println("Server started at port 8080")
	// listens on the TCP network address and then handles requests on incoming connections.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
