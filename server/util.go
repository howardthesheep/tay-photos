package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Creates a basic string with some info about the received HTTP Request
func requestString(r *http.Request) string {
	return r.Method + " " + r.RequestURI
}

// Prints the Header & Body of a received HTTP Request
func debugRequest(r *http.Request) {
	println("Req Header\n---------")
	for name, values := range r.Header {
		for _, value := range values {
			println(name, value)
		}
	}

	println("Req Body\n--------")
	buf, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		log.Fatalf("Error reading request body: %s", bodyErr)
		return
	}
	println(buf)
}

// Given an API endpoint, this function returns the subtree
// Ex.
//		Input = '/user/login'
//		Output = '/login'
func trimParentEndpoint(str string, parentEndpoint string) string {
	var endpoint string

	urlTokens := strings.Split(str, parentEndpoint)
	if len(urlTokens) == 1 {
		endpoint = "/"
	} else {
		endpoint = urlTokens[1]
	}

	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		log.Printf("Error parsing string into URL: %s", err)
		return ""
	}

	return endpointUrl.Path
}
