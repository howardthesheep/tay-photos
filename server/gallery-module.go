package main

import (
	"log"
	"net/http"
)

// API Module which handles all /gallery subtree endpoints
func galleryModule(w http.ResponseWriter, r *http.Request) {
	println("Gallery Module Request: " + requestString(r))
	endpoint := trimParentEndpoint(r.RequestURI, "/gallery")

	switch endpoint {
	case "/":
		galleryCrud(w, r)
		break
	default:
		println("Unrecognized Endpoint: " + endpoint)
	}
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
