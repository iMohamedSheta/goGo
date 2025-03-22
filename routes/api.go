package routes

import (
	"imohamedsheta/gocrud/controllers"

	"github.com/gorilla/mux"
)

func RegisterApiRoutes(apiRouter *mux.Router) {

	registerTodoRoutes(apiRouter)
}

func registerTodoRoutes(apiRouter *mux.Router) {
	// Todos example of restful endpoints
	todoController := controllers.TodoController{}

	apiRouter.HandleFunc("/todos", todoController.Index).Methods("GET")
	apiRouter.HandleFunc("/todos/{id}", todoController.Show).Methods("GET")

	apiRouter.HandleFunc("/todos", todoController.Create).Methods("POST")
	apiRouter.HandleFunc("/todos/{id}", todoController.Update).Methods("PUT")

	apiRouter.HandleFunc("/todos/{id}", todoController.Delete).Methods("DELETE")
}
