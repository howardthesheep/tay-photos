package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Gallery struct {
	Id           string    `json:"id"`
	Owner        string    `json:"owner"`
	Name         string    `json:"name"`
	CanonicalUrl string    `json:"canonical_url"`
	CreateTime   time.Time `json:"create_time"`
}

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
	urlObj, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Printf("Error during url parsing: %s", urlObj.String())
		w.WriteHeader(500)
		return
	}

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
		err = getGallery(w, urlObj.Query().Get("id"))
		break
	default:
		println("Unhandled Method on Photo: " + r.Method)
	}

	if err != nil {
		log.Printf("Error during galleryCrud: %s", err)
		w.WriteHeader(500)
		return
	}
}

// TODO: Implement these
func createGallery() {}
func deleteGallery() {}
func updateGallery() {}
func getGallery(w http.ResponseWriter, id string) error {
	// Query the db
	db := GetDatabase()
	rows := db.stmts["getGallery"].QueryRow(id)

	// Scan rows into struct
	gallery := Gallery{}
	err := rows.Scan(&gallery.Id, &gallery.Owner, &gallery.Name, &gallery.CanonicalUrl, &gallery.CreateTime)
	if err != nil {
		return err
	}

	// Parse struct into json
	galleryJson, err := json.Marshal(gallery)
	if err != nil {
		return err
	}

	// Write json to client
	_, err = w.Write(galleryJson)
	if err != nil {
		return err
	}

	return nil
}

func getGalleryPhotos(w http.ResponseWriter, r *http.Request) {

}
