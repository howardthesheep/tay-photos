package main

import (
	"log"
	"net/http"
)

/***************************************
 *		Routing Data Structures
 **************************************/

type RouteHandler struct {
	baseHandler http.HandlerFunc
	subHandlers []RouteHandler
}

type Route struct {
	base      string  // The 'root' of the route tree (Ex. /photos)
	subRoutes []Route // The subchildren of the route root
}

type Router struct {
	routes        []Route        // Represent the universe of endpoints in the application
	routeHandlers []RouteHandler // Represents the handlers of above endpoints
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

	r.routes = []Route{
		Route{
			"/",
			[]Route{
				Route{
					"/photo/",
					nil,
				},
				Route{
					"/user/",
					nil,
				},
				Route{
					"/gallery/",
					nil,
				},
			},
		},
	}

	r.routeHandlers = []RouteHandler{
		RouteHandler{
			fileServerMiddleware(fs),
			[]RouteHandler{
				{
					photoModule,
					nil,
				},
				{
					userModule,
					nil,
				},
				{
					galleryModule,
					nil,
				},
			},
		},
	}

	r._applyRouteHandlers()
}

/***************************************
 *			Helper Functions
 **************************************/

func (r Router) _applyRouteHandlers() {
	for i, route := range r.routes {
		routeHandler := r.routeHandlers[i]

		http.HandleFunc(route.base, routeHandler.baseHandler)
		_applySubrouteHandlers(route.subRoutes, routeHandler.subHandlers)
	}
}

// Recursively applies all the subroute handlers to the proper subroute
func _applySubrouteHandlers(subRoute []Route, subHandlers []RouteHandler) {
	if subRoute == nil {
		return
	}

	for i, route := range subRoute {
		handler := subHandlers[i]

		http.HandleFunc(route.base, handler.baseHandler)
		_applySubrouteHandlers(route.subRoutes, handler.subHandlers)
	}
}
