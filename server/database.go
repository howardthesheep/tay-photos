package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB
var dbStmts = make(map[QueryName]*sql.Stmt)

var preparedStmts = map[QueryName]string{
	CreateUser:                  "INSERT INTO Users (id, name, username, email, password, apiToken) VALUES (?,?,?,?,?,?);",
	DeleteUser:                  "DELETE FROM Users WHERE id=?;",
	UpdateUser:                  "UPDATE Users SET name=?, username=?, email=? WHERE apiToken=?",
	GetUser:                     "SELECT name,username,email FROM Users WHERE id=?;",
	UserLogin:                   "SELECT apiToken FROM Users WHERE username=? AND password=?",
	EmailLogin:                  "SELECT apiToken FROM Users WHERE email=? AND password=?",
	GetUserByEmail:              "SELECT * FROM Users WHERE email=?;",
	GetGallery:                  "SELECT * FROM Galleries WHERE id=?;",
	GetGalleryPhotosByGalleryId: "SELECT gallery,collection,id FROM GalleryPhotos WHERE gallery=?;",
	GetGalleryPhotoFilePath:     "SELECT photoPath FROM GalleryPhotos WHERE id=?;",
}

type QueryName int

const (
	CreateUser QueryName = iota
	DeleteUser
	UpdateUser
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
	// Read the .env file storing credentials
	envFile, err := os.Open("../.env")
	if err != nil {
		log.Printf("Error opening file" + err.Error())
		return ""
	}

	// Read contents into variable
	var contents []byte
	contents, err = io.ReadAll(envFile)
	if err != nil {
		log.Printf("Error reading file" + err.Error())
		return ""
	}

	// Parse key value pairs from .env file
	var keyVal = map[string]string{}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		pair := strings.Split(line, "=")
		keyVal[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		keyVal["username"],
		keyVal["password"],
		keyVal["hostname"],
		keyVal["dbname"],
	)
}
