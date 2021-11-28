package main

import (
	"log"
	"net/http"
)

// API Module which handles all /user subtree endpoints
func userModule(w http.ResponseWriter, r *http.Request) {
	println("User Module Request: " + requestString(r))
	endpoint := trimParentEndpoint(r.RequestURI, "/user")

	switch endpoint {
	case "/":
		userCrud(w, r)
		break
	case "/login":
		login(w, r)
		break
	default:
		println("Unrecognized Endpoint: " + endpoint)
	}
}

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

// Authenticates a user and then returns an apiToken for privileged actions
func login(w http.ResponseWriter, r *http.Request) {

}
