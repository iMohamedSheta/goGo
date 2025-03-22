package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ListAllRoutes(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routes []map[string]interface{}

		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()
			if err == nil {
				methods, _ := route.GetMethods()

				routes = append(routes, map[string]interface{}{
					"path":    path,
					"methods": methods,
				})
			}
			return nil
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes)
	}
}
