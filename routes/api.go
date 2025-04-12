package routes

import (
	"imohamedsheta/gocrud/app/controllers"

	"github.com/gorilla/mux"
)

func RegisterApiRoutes(apiRouter *mux.Router) {

	registerTodoRoutes(apiRouter)
}

func registerTodoRoutes(apiRouter *mux.Router) {
	// Todos example of restful endpoints
	todoController := controllers.TodoController{}

	// swagger:route GET /todos todos listTodos
	apiRouter.HandleFunc("/todos", todoController.Index).Methods("GET")
	// swagger:route GET /todos/{id} todos showTodo
	apiRouter.HandleFunc("/todos/{id}", todoController.Show).Methods("GET")

	// swagger:route POST /todos todos createTodo
	// Create a new todo
	// responses:
	// 	201: Todo
	// 422: ValidationError
	// 500: InternalServerError
	apiRouter.HandleFunc("/todos", todoController.Create).Methods("POST")
	// swagger:route PUT /todos/{id} todos updateTodo
	// Update a todo
	// responses:
	// 200: Todo
	// 422: ValidationError
	// 500: InternalServerError
	apiRouter.HandleFunc("/todos/{id}", todoController.Update).Methods("PUT")

	// swagger:route DELETE /todos/{id} todos deleteTodo
	// Delete a todo
	// responses:
	// 200: Todo
	// 500: InternalServerError
	// 422: ValidationError
	// 404: NotFound
	// 401: Unauthorized
	apiRouter.HandleFunc("/todos/{id}", todoController.Delete).Methods("DELETE")
}
