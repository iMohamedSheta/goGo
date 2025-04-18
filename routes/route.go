package routes

import (
	"imohamedsheta/gocrud/helpers"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()
	mainRoutes(router)

	apiRouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	RegisterApiRoutes(apiRouter)

	return router
}

func mainRoutes(router *mux.Router) {

	// Register the "list all routes" endpoint
	router.HandleFunc("/routes", helpers.ListAllRoutes(router)).Methods(http.MethodGet)
}
