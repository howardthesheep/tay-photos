package main

import (
	"net/http"
)

func main() {
	router := Router{}
	router.initWebsiteRoutes()
	router.startWebsite()
}

// TODO: move this to middleware file
// Extra logic on top of the /www file server to ensure authentication, logging, etc.
func fileServerMiddleware(fs http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		println(requestString(request))

		// TODO: Ensure authentication for access to gallery editor pages & stuff

		fs.ServeHTTP(writer, request)
	}
}
