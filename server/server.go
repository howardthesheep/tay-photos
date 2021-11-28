package main

import (
	"log"
	"net/http"
)

func main() {
	apiLocation := "http://localhost:6969"

	// Expose /www files
	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", fs)

	// Expose API Modules
	http.HandleFunc("/photo/", photoModule)
	http.HandleFunc("/user/", userModule)
	http.HandleFunc("/gallery/", galleryModule)

	println("Running webserver at " + apiLocation)
	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		log.Fatalf(`Error during http server start: $err`)
		return
	}
}
