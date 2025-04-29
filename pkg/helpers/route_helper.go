package helpers

import (
	"encoding/json"
	"imohamedsheta/gocrud/pkg/router"
	"net/http"

	"github.com/gorilla/mux"
)

// func ListAllRoutes(router *mux.Router) http.HandlerFunc {
func ListAllRoutes(router *router.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routes []map[string]interface{}

		// Walk through all routes in the main router
		router.Walk(func(route *mux.Route, muxRouter *mux.Router, ancestors []*mux.Route) error {
			// Get the route path
			path, err := route.GetPathTemplate()
			if err == nil {
				// Get the HTTP methods for this route (e.g., GET, POST, etc.)
				methods, _ := route.GetMethods()

				if len(methods) == 1 {
					routes = append(routes, map[string]interface{}{
						"path":    path,
						"methods": methods[0],
					})
				} else if len(methods) > 1 {
					routes = append(routes, map[string]interface{}{
						"path":    path,
						"methods": methods,
					})
				}

			}
			return nil
		})

		// Output the collected routes in JSON format
		w.Header().Set("Content-Type", "application/json")
		// Output middleware functions as their name (if needed)
		json.NewEncoder(w).Encode(routes)
	}
}
