package main

import (
	"log"
	"net/http"
)

/***************************************
 *		Routing Data Structures
 **************************************/

type Route struct {
	base        string // The 'root' of the route tree (Ex. /photos)
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
	fs := http.FileServer(http.Dir("./www"))

	r.routes = Route{
		"/",
		fileServerMiddleware(fs),
		[]Route{
			Route{
				"/photo/",
				photoModule,
				nil,
			},
			Route{
				"/user/",
				userModule,
				nil,
			},
			Route{
				"/gallery/",
				galleryModule,
				nil,
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
	// Register the base route
	route := r.routes
	http.HandleFunc(route.base, route.baseHandler)

	// Apply all the subroutes
	_applySubrouteHandlers(r.routes.subRoutes)
}

// Recursively applies all the subroutes and their associated handlers
func _applySubrouteHandlers(subRoutes []Route) {
	if subRoutes == nil {
		return
	}

	for _, route := range subRoutes {
		http.HandleFunc(route.base, route.baseHandler)
		_applySubrouteHandlers(route.subRoutes)
	}
}
