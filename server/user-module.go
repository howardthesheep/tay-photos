package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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
	userCrud(w, r)
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
	res, err := db.stmts[CreateUser].Exec(userId.String(), userData.Name, userData.Username, userData.Email, string(hash), apiToken.String())
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
	res, err := db.stmts[DeleteUser].Exec(userData.Id)
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
	res, err := db.stmts[UpdateUser].Exec(userData.Name, userData.Username, userData.Email, userData.ApiToken)
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
	row := db.stmts[GetUser].QueryRow(id)
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

	// Reject request if not POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse POST body contents
	var bodyBytes []byte
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body bytes: %s", err)
		return
	}

	log.Printf("Recieved login data: %s", string(bodyBytes))

	// Unmarshal user login data from JSON
	userData := UserLoginData{}
	err = json.Unmarshal(bodyBytes, &userData)
	if err != nil {
		log.Printf("Error unmarshaling body: %s", err)
		w.WriteHeader(400)
		return
	}

	// Hash users password
	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		return
	}

	// Ask the db for user associated with email/username & password hash combination
	db := GetDatabase()
	var row *sql.Row
	if strings.Contains(userData.Username, "@") {
		row = db.stmts[EmailLogin].QueryRow(userData.Username, hash)
	} else {
		row = db.stmts[UserLogin].QueryRow(userData.Username, hash)
	}

	// Return a users stored apiToken to pass off to them to use for authenticating future requests
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

	// If apiToken is empty, generate them one
	if apiToken == "" {
		token := jwt.New(jwt.SigningMethodEdDSA)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(256 * time.Hour)
		claims["authorized"] = true
		claims["username"] = userData.Username

		// Sign our token with secret JWTKey from .env file
		env := GetEnv()
		apiToken, err = token.SignedString(env["JWTKey"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error signing api token: %s", err)
			return
		}

		// Update the api token associated with our user in the database
		_, err = db.stmts[UpdateUserApiToken].Exec(apiToken, userData.Username, hash)
		if err != nil {
			log.Printf("Error updating user api token in database: %s", err)
			return
		}
	}

	// Stuff apiToken into map to be marshalled into json and sent back to client
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
