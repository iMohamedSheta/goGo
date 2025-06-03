package router

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Define a new custom type that embeds mux.Router
type Router struct {
	*mux.Router
}

var (
	router     *Router
	routerOnce sync.Once
)

// Instance returns the singleton router instance
// It is thread-safe and will only initialize the router once
func Instance() *Router {
	routerOnce.Do(func() {
		router = &Router{
			Router: mux.NewRouter(),
		}
	})
	return router
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
		Router: groupRouter,
	}

	// Pass the group router to the function to register routes
	f(group)
}
