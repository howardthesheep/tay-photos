package main

import (
	"errors"
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

		page := request.RequestURI
		switch page {
		case "/dashboard.html":
			fallthrough
		case "/gallery.html":
			err := authenticateIfNecessary(request)
			if err != nil {
				switch err.Error() {
				case "unauthenticated":
					http.Redirect(writer, request, "/login.html", 301)
					return
				}
			}
			fallthrough
		default:
			fs.ServeHTTP(writer, request)

		}
	}
}

func authenticateIfNecessary(request *http.Request) error {
	authToken := request.Header.Get("Authorization")
	if len(authToken) < 5 {
		return errors.New("unauthenticated")
	}

	return nil
}
