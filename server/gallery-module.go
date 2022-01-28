package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Gallery struct {
	Id           string    `json:"id"`
	Owner        string    `json:"owner"`
	Name         string    `json:"name"`
	CanonicalUrl string    `json:"canonical_url"`
	CreateTime   time.Time `json:"create_time"`
}

type GalleryPhoto struct {
	Id         string `json:"id"`
	Gallery    string `json:"gallery"`
	Collection string `json:"collection"`
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

// Gets the json representation of all the GalleryPhotos rows
// associated with the provided Gallery.Id in the query args
func getGalleryPhotos(w http.ResponseWriter, r *http.Request) {
	println("Getting gallery photos...")
	urlObj, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Printf("Error during url parsing: %s", urlObj)
		w.WriteHeader(500)
		return
	}

	id := urlObj.Query().Get("id")
	db := GetDatabase()

	rows, err := db.stmts["getGalleryPhotosByGalleryId"].Query(id)
	if err != nil {
		log.Printf("Error running db query: %s", err)
		w.WriteHeader(500)
		return
	}

	// Scan rows into struct
	var galleryPhotoData []GalleryPhoto
	for rows.Next() {
		photoData := GalleryPhoto{}
		err = rows.Scan(&photoData.Gallery, &photoData.Collection, &photoData.Id)
		if err != nil {
			log.Printf("Error row to struct scanning: %s", err)
			w.WriteHeader(500)
			return
		}

		galleryPhotoData = append(galleryPhotoData, photoData)
	}

	// Parse struct into json
	galleryJson, err := json.Marshal(galleryPhotoData)
	if err != nil {
		log.Printf("Error marshalling json into struct: %s", err)
		w.WriteHeader(500)
		return
	}

	// Write json to client
	_, err = w.Write(galleryJson)
	if err != nil {
		log.Printf("Error during writing json to client: %s", err)
		w.WriteHeader(500)
		return
	}
}

// Sends the serialized photo based on the id provided in the query args
func getGalleryPhoto(w http.ResponseWriter, r *http.Request) {
	urlObj, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Printf("Error during url parsing: %s", urlObj)
		w.WriteHeader(500)
		return
	}
	id := urlObj.Query().Get("id")

	// Crack open DB and run query
	db := GetDatabase()
	rows, err := db.stmts["getGalleryPhotoFilePath"].Query(id)
	if err != nil {
		log.Printf("Error querying db: %s", err)
		w.WriteHeader(500)
		return
	}

	var photoPath string
	for rows.Next() {
		err = rows.Scan(&photoPath)
		if err != nil {
			log.Printf("Error row to struct scanning: %s", err)
			w.WriteHeader(500)
			return
		}
	}

	// Read file off of file system
	fileBytes, err := os.ReadFile("/home/pjones" + photoPath)
	if err != nil {
		log.Printf("Error during photo open on filesystem: %s", err)
		w.WriteHeader(500)
		return
	}

	// Write file in http response back to client
	w.Header().Set("Content-Type", "image/*")
	w.Header().Set("Content-Disposition", "inline;")
	_, err = w.Write(fileBytes)
	if err != nil {
		log.Printf("Error writing file bytes to client: %s", err)
		w.WriteHeader(500)
		return
	}
}
