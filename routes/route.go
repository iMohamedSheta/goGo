package routes

import (
	"imohamedsheta/gocrud/pkg/helpers"
	"imohamedsheta/gocrud/pkg/router"
	"imohamedsheta/gocrud/pkg/support"
)

// RegisterRoutes registers all routes for the application.
func RegisterRoutes() *router.Router {
	r := router.Instance()

	mainRoutes(r)

	r.Group("api", "/api/v1", nil, func(apiRouter *router.Router) {
		RegisterApiRoutes(apiRouter)
	})

	return r
}

func mainRoutes(r *router.Router) {

	// Register the "list all routes" endpoint
	if support.IsDev() {
		r.Get("/routes", helpers.ListAllRoutes(router.Instance())).Name("list_all_routes")
	}
}
