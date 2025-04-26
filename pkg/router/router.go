package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Define a new custom type that embeds mux.Router
type Router struct {
	*mux.Router
	routeMiddleware map[string][]mux.MiddlewareFunc
	groups          map[string]*Router // <-- track groups
}

// Create a new instance of Router
func New() *Router {
	return &Router{
		Router:          mux.NewRouter(),
		routeMiddleware: make(map[string][]mux.MiddlewareFunc),
		groups:          make(map[string]*Router),
	}
}

// Define a custom Get method on the Router type
func (r *Router) Get(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("GET")
}

// Define a custom Post method on the Router type
func (r *Router) Post(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("POST")
}

// Define a custom Put method on the Router type
func (r *Router) Put(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("PUT")
}

// Define a custom Delete method on the Router type
func (r *Router) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("DELETE")
}

// Define a custom Options method on the Router type
func (r *Router) Options(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("OPTIONS")
}

// Define a custom Patch method on the Router type
func (r *Router) Patch(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("PATCH")
}

// Group method to create a subrouter with custom middlewares
func (r *Router) Group(name string, path string, middleware []mux.MiddlewareFunc, f func(r *Router)) {
	groupRouter := r.PathPrefix(path).Subrouter()

	// Apply middleware to the group subrouter
	for _, mw := range middleware {
		groupRouter.Use(mw)
	}

	// Create a new Router instance for the group, tracking middleware
	group := &Router{
		Router:          groupRouter,
		routeMiddleware: make(map[string][]mux.MiddlewareFunc),
		groups:          make(map[string]*Router),
	}

	r.routeMiddleware[name] = middleware
	r.groups[name] = group

	// Pass the group router to the function to register routes
	f(group)
}

func (r *Router) GetGroup(name string) *Router {
	return r.groups[name]
}

// Method to retrieve middlewares applied to a route
func (r *Router) GetMiddlewaresForRoute(path string) []mux.MiddlewareFunc {
	return r.routeMiddleware[path]
}

func (r *Router) GetMiddlewaresWithGroup() map[string][]mux.MiddlewareFunc {
	return r.routeMiddleware
}
