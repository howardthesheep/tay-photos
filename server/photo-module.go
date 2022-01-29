package main

import (
	"log"
	"net/http"
)

// API Module which handles all /photo subtree endpoints
func photoModule(w http.ResponseWriter, r *http.Request) {
	println("Photo Module Request: " + requestString(r))
	photoCrud(w, r)
}

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
