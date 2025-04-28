package routes

import (
	"imohamedsheta/gocrud/pkg/helpers"
	"imohamedsheta/gocrud/pkg/router"
	"imohamedsheta/gocrud/pkg/support"
)

// func RegisterRoutes() *mux.Router {
// 	router := mux.NewRouter()
// 	mainRoutes(router)

// 	apiRouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
// 	RegisterApiRoutes(apiRouter)

// 	return router
// }

// func mainRoutes(router *mux.Router) {

// 	// Register the "list all routes" endpoint
// 	router.HandleFunc("/routes", helpers.ListAllRoutes(router)).Methods("GET")
// }

func RegisterRoutes() *router.Router {
	r := router.Instance()

	mainRoutes(r)

	r.Group("api", "/api/v1", nil, func(apiRouter *router.Router) {
		RegisterApiRoutes(apiRouter)
	})

	// apiRouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	return r
}

func mainRoutes(r *router.Router) {

	// Register the "list all routes" endpoint
	if support.IsDev() {
		r.Get("/routes", helpers.ListAllRoutes(router.Instance())).Name("listAllRoutes")
	}
}
