package routes

import (
	"github.com/iMohamedSheta/xapp/pkg/helpers"
	"github.com/iMohamedSheta/xapp/pkg/router"
	"github.com/iMohamedSheta/xapp/pkg/support"
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
