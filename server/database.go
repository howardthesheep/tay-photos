package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// QueryName is a custom type which takes advantage of IDE autofill for seeing the universe
// of sql.Stmt held in the dbStmts map.
type QueryName int

const (
	CreateUser QueryName = iota
	DeleteUser
	UpdateUser
	UpdateUserApiToken
	GetUser
	GetUserByEmail
	UserLogin
	EmailLogin
	GetGallery
	GetGalleryPhotosByGalleryId
	GetGalleryPhotoFilePath
)

// DBConnection holds our active db connection and access to prepared queries
type DBConnection struct {
	database *sql.DB
	stmts    map[QueryName]*sql.Stmt
}

// Globals
var database *sql.DB                        // Connector to DB
var dbStmts = make(map[QueryName]*sql.Stmt) // Map of queries available to the programmer using this

var preparedStmts = map[QueryName]string{
	CreateUser:                  "INSERT INTO Users (id, name, username, email, password, apiToken) VALUES (?,?,?,?,?,?);",
	DeleteUser:                  "DELETE FROM Users WHERE id=?;",
	UpdateUser:                  "UPDATE Users SET name=?, username=?, email=? WHERE apiToken=?",
	UpdateUserApiToken:          "Update Users SET apiToken=? WHERE username=? AND password=?",
	GetUser:                     "SELECT name,username,email FROM Users WHERE id=?;",
	UserLogin:                   "SELECT apiToken FROM Users WHERE username=? AND password=?",
	EmailLogin:                  "SELECT apiToken FROM Users WHERE email=? AND password=?",
	GetUserByEmail:              "SELECT * FROM Users WHERE email=?;",
	GetGallery:                  "SELECT * FROM Galleries WHERE id=?;",
	GetGalleryPhotosByGalleryId: "SELECT gallery,collection,id FROM GalleryPhotos WHERE gallery=?;",
	GetGalleryPhotoFilePath:     "SELECT photoPath FROM GalleryPhotos WHERE id=?;",
}

// GetDatabase returns the object for running SQL queries on Golang Backend
func GetDatabase() DBConnection {
	if database != nil && dbStmts != nil {
		return DBConnection{
			database: database,
			stmts:    dbStmts,
		}
	}

	openConnection()
	prepareStatements()
	return DBConnection{
		database: database,
		stmts:    dbStmts,
	}
}

// openConnection opens a connection to backend DB and ensures it can be reached
func openConnection() {
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		log.Fatalf(`Error when getting DB connection: %s`, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Error when pinging database: %s", err)
		return
	}

	database = db
}

// prepareStatements adds our prepared SQL queries to the DB connection
func prepareStatements() {
	for key, value := range preparedStmts {
		prepedStmt, err := database.Prepare(value)
		if err != nil {
			log.Fatalf("Error preparing query: %s", err)
			return
		}

		dbStmts[key] = prepedStmt
	}
}

// closeConnection closes our connection to backend DB
func closeConnection() {
	if database != nil {
		err := database.Close()
		if err != nil {
			log.Fatalf("Error closing db connection: %s", err)
			return
		}
	}
}

// Helper function for converting db connection info into usable form
func dsn() string {
	env := GetEnv()

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		env["username"],
		env["password"],
		env["hostname"],
		env["dbname"],
	)
}
