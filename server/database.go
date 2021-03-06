package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "tayphoto_user"
	password = "tayphoto_pass"
	hostname = "127.0.0.1:3306"
	dbname   = "tay_photos"
)

var database *sql.DB
var dbStmts = make(map[string]*sql.Stmt)

var preparedStmts = map[string]string{
	"createUser":                  "INSERT INTO Users (id, name, username, email, password, apiToken) VALUES (?,?,?,?,?,?);",
	"deleteUser":                  "DELETE FROM Users WHERE id=?;",
	"updateUser":                  "UPDATE Users SET name=?, username=?, email=? WHERE apiToken=?",
	"getUser":                     "SELECT name,username,email FROM Users WHERE id=?;",
	"userLogin":                   "SELECT apiToken FROM Users WHERE username=? AND password=?",
	"emailLogin":                  "SELECT apiToken FROM Users WHERE email=? AND password=?",
	"getUserByEmail":              "SELECT * FROM Users WHERE email=?;",
	"getGallery":                  "SELECT * FROM Galleries WHERE id=?;",
	"getGalleryPhotosByGalleryId": "SELECT gallery,collection,id FROM GalleryPhotos WHERE gallery=?;",
	"getGalleryPhotoFilePath":     "SELECT photoPath FROM GalleryPhotos WHERE id=?;",
}

// DBConnection holds our active db connection and access to prepared queries
type DBConnection struct {
	database *sql.DB
	stmts    map[string]*sql.Stmt
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
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbname)
}
