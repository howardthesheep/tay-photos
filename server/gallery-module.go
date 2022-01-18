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

// API Module which handles all /gallery.html subtree endpoints
func galleryModule(w http.ResponseWriter, r *http.Request) {
	println("Gallery Module Request: " + requestString(r))
	galleryCrud(w, r)
}

// Handles Requests for CRUD operations on Galleries
func galleryCrud(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.Method {
	case "POST":
		createGallery(w, r)
		break
	case "DELETE":
		deleteGallery(w, r)
		break
	case "PUT":
		updateGallery(w, r)
		break
	case "GET":
		err = getGallery(w, r)
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
func createGallery(w http.ResponseWriter, r *http.Request) {}
func deleteGallery(w http.ResponseWriter, r *http.Request) {}
func updateGallery(w http.ResponseWriter, r *http.Request) {}
func getGallery(w http.ResponseWriter, r *http.Request) error {
	urlObj, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Printf("Error during url parsing: %s", urlObj)
		w.WriteHeader(500)
		return err
	}

	id := urlObj.Query().Get("id")

	// Query the db
	db := GetDatabase()
	rows := db.stmts["getGallery"].QueryRow(id)

	// Scan rows into struct
	gallery := Gallery{}
	err = rows.Scan(&gallery.Id, &gallery.Owner, &gallery.Name, &gallery.CanonicalUrl, &gallery.CreateTime)
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
	println("Getting gallery.html photos...")
}
