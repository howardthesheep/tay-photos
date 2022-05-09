package main

import (
	"log"
	"net/http"
	"strings"
)

/***************************************
 *		Routing Data Structures
 **************************************/

type Route struct {
	name        string // The 'root' of the route tree (Ex. /photos)
	baseHandler http.HandlerFunc
	subRoutes   []Route // The subchildren of the route root
}

type Router struct {
	routes Route // Represent the universe of endpoints in the application
}

/***************************************
 *			Router Functions
 **************************************/

func (r Router) startWebsite() {
	apiLocation := "http://localhost:6969"

	println("Running webserver at " + apiLocation)
	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		log.Fatalf(`Error during http server start: $err`)
		return
	}
}

func (r Router) initWebsiteRoutes() {
	r.routes = Route{
		"/",
		http.NotFound,
		[]Route{
			Route{
				"user",
				userModule,
				[]Route{
					Route{
						"login",
						login,
						nil,
					},
				},
			},
			Route{
				"gallery",
				galleryModule,
				[]Route{
					Route{
						"photos",
						getGalleryPhotos,
						nil,
					},
					Route{
						"photo",
						getGalleryPhoto,
						nil,
					},
				},
			},
		},
	}

	r._applyRouteHandlers()
}

/***************************************
 *			Helper Functions
 **************************************/

// Registers all of a routers routes and associated handlers with HTTP
func (r Router) _applyRouteHandlers() {
	// Register the name route
	route := r.routes
	http.HandleFunc(route.name, route.baseHandler)

	// Apply all the subroutes
	_applySubrouteHandlers(route.subRoutes, route.name)
}

// Recursively applies all the subroutes and their associated handlers
func _applySubrouteHandlers(subRoutes []Route, baseName string) {
	if subRoutes == nil {
		return
	}

	for _, route := range subRoutes {
		var fullPath string
		if strings.HasSuffix(baseName, "/") || strings.HasPrefix(route.name, "/") {
			fullPath = baseName + route.name
		} else {
			fullPath = baseName + "/" + route.name
		}

		http.HandleFunc(fullPath, route.baseHandler)
		_applySubrouteHandlers(route.subRoutes, fullPath)
	}
}
