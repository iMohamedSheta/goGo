package helpers

import (
	"encoding/json"
	"imohamedsheta/gocrud/pkg/router"
	"net/http"
	"strings"

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

				// Look up the middleware applied to this specific route
				routeMiddlewares := router.GetMiddlewaresForRoute(path)

				// Check for group-level middleware
				groupMiddleware := make([]mux.MiddlewareFunc, 0)
				for groupName, groupMw := range router.GetMiddlewaresWithGroup() {
					// Check if the route path starts with the group path prefix
					groupPrefixPath := "/" + groupName // e.g., "api", so the prefix is "/api"
					if strings.HasPrefix(path, groupPrefixPath) {
						// If the route matches the group, add its middleware
						groupMiddleware = append(groupMiddleware, groupMw...)
					}
				}

				// Combine route-specific and group-specific middlewares
				allMiddlewares := append(routeMiddlewares, groupMiddleware...)

				// Store the route information (path, methods, and middlewares)
				routes = append(routes, map[string]interface{}{
					"path":        path,
					"methods":     methods,
					"middlewares": allMiddlewares, // This now holds the actual middleware functions
				})
			}
			return nil
		})

		// Output the collected routes in JSON format
		w.Header().Set("Content-Type", "application/json")
		// Output middleware functions as their name (if needed)
		json.NewEncoder(w).Encode(routes)
	}
}
