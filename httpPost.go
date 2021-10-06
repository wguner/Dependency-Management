package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// A map is serialized into JSON string with json.Marshal. 
	values := map[string]string{"name": "John Doe", "occupation": "gardener"}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	// sending a post request to httpbin
	// data is sent with PostForm
	// PostForm issues a POST to the specified URL, with data's keys and values URL-encoded as the request body. 
	resp, err := http.Post("https://httpbin.org/post", "application/json",
		bytes.NewBuffer(json_data))
	// When we post the data, we set the content type to application/json. 
	if err != nil {
		log.Fatal(err)
	}

	// decoding postform request
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	// printing to console
	fmt.Println(res["json"])
}
