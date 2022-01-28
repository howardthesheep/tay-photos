package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

type UserLoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authenticates a user and then returns an apiToken for privileged actions
func login(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body bytes: %s", err)
		return
	}

	log.Printf("Recieved login data: %s", string(bodyBytes))

	userData := UserLoginData{}
	err = json.Unmarshal(bodyBytes, &userData)
	if err != nil {
		log.Printf("Error unmarshaling body: %s", err)
		return
	}

	// Hash users password
	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		return
	}

	db := GetDatabase()
	var row *sql.Row
	if strings.Contains(userData.Username, "@") {
		row = db.stmts["emailLogin"].QueryRow(userData.Username, hash)
	} else {
		row = db.stmts["userLogin"].QueryRow(userData.Username, hash)
	}

	var apiToken string
	err = row.Scan(&apiToken)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(401)
			_, err = w.Write([]byte("Invalid username/password combination"))
			if err != nil {
				log.Printf("Error writing 401 to client: %s", err)
				return
			}
		} else {
			log.Printf("Error scanning apiToken from db: %s", err)
			return
		}
	}

	// Stuff apiToken into map to be marshalled into json
	jsonMap := make(map[string]string)
	jsonMap["apiToken"] = apiToken
	jsonBytes, err := json.Marshal(jsonMap)

	// Write json response back to client
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Error sending response to client: %s", err)
		return
	}
}
