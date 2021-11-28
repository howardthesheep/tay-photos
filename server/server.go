package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	apiLocation := "http://localhost:6969"

	// Expose /www files
	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", fs)

	// Expose API Modules
	http.HandleFunc("/photo/", photoModule)
	http.HandleFunc("/user/", userModule)
	http.HandleFunc("/gallery/", galleryModule)

	println("Running webserver at " + apiLocation)
	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		log.Fatalf(`Error during http server start: $err`)
		return
	}
}

/***********************************
 *
 * 			API Modules
 *
***********************************/

// API Module which handles all /photo subtree endpoints
func photoModule(w http.ResponseWriter, r *http.Request) {
	println("Photo Module Request: " + requestString(r))
	var endpoint string

	urlTokens := strings.Split(r.RequestURI, "/photo")
	if len(urlTokens) == 1 {
		endpoint = "/"
	} else {
		endpoint = urlTokens[1]
	}

	switch endpoint {
	case "/":
		photoCrud(w, r)
		break
	default:
		println("Unrecognized Endpoint: " + endpoint)
	}
}

// API Module which handles all /user subtree endpoints
func userModule(w http.ResponseWriter, r *http.Request) {
	println("User Module Request: " + requestString(r))
	var endpoint string

	urlTokens := strings.Split(r.RequestURI, "/user")
	if len(urlTokens) == 1 {
		endpoint = "/"
	} else {
		endpoint = urlTokens[1]
	}

	switch endpoint {
	case "/":
		userCrud(w, r)
		break
	default:
		println("Unrecognized Endpoint: " + endpoint)
	}
}

// API Module which handles all /gallery subtree endpoints
func galleryModule(w http.ResponseWriter, r *http.Request) {
	println("Gallery Module Request: " + requestString(r))
	var endpoint string

	urlTokens := strings.Split(r.RequestURI, "/gallery")
	if len(urlTokens) == 1 {
		endpoint = "/"
	} else {
		endpoint = urlTokens[1]
	}

	switch endpoint {
	case "/":
		galleryCrud(w, r)
		break
	default:
		println("Unrecognized Endpoint: " + endpoint)
	}
}

/***********************************
 *
 * 		   Helper Functions
 *
***********************************/

// Handles Requests for CRUD operations on Photos
func photoCrud(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createPhoto()
		break
	case "DELETE":
		deletePhoto()
		break
	case "PUT":
		updatePhoto()
		break
	case "GET":
		getPhoto()
		break
	default:
		log.Fatalf("Unhandled Method on Photo: " + r.Method)
	}
}

// TODO: Implement these
func createPhoto() {}
func deletePhoto() {}
func updatePhoto() {}
func getPhoto()    {}

// Handles Requests for CRUD operations on Users
func userCrud(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createUser()
		break
	case "DELETE":
		deleteUser()
		break
	case "PUT":
		updateUser()
		break
	case "GET":
		getUser()
		break
	default:
		log.Fatalf("Unhandled Method on Photo: " + r.Method)
	}
}

// TODO: Implement these
func createUser() {
	println("Creating User...")
}
func deleteUser() {
	println("Deleting User...")
}
func updateUser() {
	println("Updating User...")
}
func getUser() {
	println("Getting User...")
}

// Handles Requests for CRUD operations on Galleries
func galleryCrud(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createGallery()
		break
	case "DELETE":
		deleteGallery()
		break
	case "PUT":
		updateGallery()
		break
	case "GET":
		getGallery()
		break
	default:
		log.Fatalf("Unhandled Method on Photo: " + r.Method)
	}
}

// TODO: Implement these
func createGallery() {}
func deleteGallery() {}
func updateGallery() {}
func getGallery()    {}

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
