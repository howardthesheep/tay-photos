package main

import (
	"log"
	"net/http"
)

func main() {
	apiLocation := "http://localhost:6969"

	// Expose /www files
	fs := http.FileServer(http.Dir("../www"))
	http.Handle("/", fileServerMiddleware(fs))

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

// Extra logic on top of the /www file server to ensure authentication, logging, etc.
func fileServerMiddleware(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println(requestString(request))

		// TODO: Ensure authentication for access to gallery editor pages & stuff

		fs.ServeHTTP(writer, request)
	})
}
