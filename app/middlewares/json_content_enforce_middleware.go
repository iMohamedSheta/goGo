package middlewares

import (
	"imohamedsheta/gocrud/pkg/response"
	"net/http"
	"strings"
)

func JSONContentTypeMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")

			// Check if the Content-Type is application/json
			if !strings.HasPrefix(contentType, "application/json") {
				response.ErrorJson(w, "Content-Type must be application/json", "invalid_content_type", http.StatusUnsupportedMediaType)
				return
			}

			// If it's valid, proceed with the request
			next.ServeHTTP(w, r)
		})
	}
}
