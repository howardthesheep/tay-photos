package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDeleteData struct {
	Id string `json:"id"`
}

type UserCreateData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	UserLoginData
}

type UserUpdateData struct {
	ApiToken string `json:"api_token"`
	UserCreateData
}

type UserBasicData struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

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
		createUser(w, r)
		break
	case "DELETE":
		deleteUser(w, r)
		break
	case "PUT":
		updateUser(w, r)
		break
	case "GET":
		getUser(w, r)
		break
	default:
		log.Fatalf("Unhandled Method on Photo: " + r.Method)
	}
}

// Creates a new User row in the database based on the info provided in the request
func createUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading client request body: %s", err)
		return
	}

	var userData UserCreateData
	err = json.Unmarshal(bodyBytes, &userData)
	if err != nil {
		log.Printf("Error unmarshalling request data: %s", err)
		return
	}

	userId := uuid.New()
	apiToken := uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %s", err)
		return
	}

	db := GetDatabase()
	res, err := db.stmts["createUser"].Exec(userId.String(), userData.Name, userData.Username, userData.Email, string(hash), apiToken.String())
	if err != nil {
		log.Printf("Error inserting new user into db: %s", err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows: %s", err)
		return
	}

	if rowsAffected == 1 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Removes a User row in the database based on the id provided in the request
func deleteUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading client request body: %s", err)
		return
	}

	var userData UserDeleteData
	err = json.Unmarshal(bodyBytes, &userData)
	if err != nil {
		log.Printf("Error parsing request json: %s", err)
		return
	}

	db := GetDatabase()
	res, err := db.stmts["deleteUser"].Exec(userData.Id)
	if err != nil {
		log.Printf("Error deleting user: %s", err)
		w.WriteHeader(500)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting amount of rows affected from delete: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Updates a user row in the database based on the info provided in the request
func updateUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		return
	}

	var userData UserUpdateData
	err = json.Unmarshal(bodyBytes, &userData)
	if err != nil {
		log.Printf("Error parsing body into json: %s", err)
		return
	}

	db := GetDatabase()
	res, err := db.stmts["updateUser"].Exec(userData.Name, userData.Username, userData.Email, userData.ApiToken)
	if err != nil {
		log.Printf("Error updating user in database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ra, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected row count: %s", err)
		return
	}

	if ra == 1 {
		w.WriteHeader(http.StatusOK)
		return
	} else if ra < 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Invalid")
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Somehow updated more than 1 row when only 1 should have updated, uh oh")
		return
	}
}

// Gets a specific user row in the database based on the provided id
func getUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("No user id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userData UserBasicData
	db := GetDatabase()
	row := db.stmts["getUser"].QueryRow(id)
	err := row.Scan(&userData.Name, &userData.Username, &userData.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user exists with provided id")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			log.Printf("Error scanning userdata from sql: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshalling struct into json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response to client: %s", err)
		return
	}
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
