package routes

import (
	"imohamedsheta/gocrud/app/controllers"
	"imohamedsheta/gocrud/app/middlewares"
	"imohamedsheta/gocrud/pkg/router"

	"github.com/gorilla/mux"
)

func RegisterApiRoutes(apiRouter *router.Router) {

	registerAuthRoutes(apiRouter)

	apiRouter.Group("auth", "/", []mux.MiddlewareFunc{middlewares.AuthMiddleware()},
		func(router *router.Router) {
			registerTodoRoutes(router)
		})

}

func registerTodoRoutes(router *router.Router) {
	// Todos example of restful endpoints
	todoController := controllers.TodoController{}

	// swagger:route GET /todos todos listTodos
	router.Get("/todos", todoController.Index).Name("ListTodos")
	// apiRouter.HandleFunc("/todos", todoController.Index).Methods("GET")
	// swagger:route GET /todos/{id} todos showTodo
	router.Get("/todos/{id}", todoController.Show).Name("ListTodos")
	// apiRouter.HandleFunc("/todos/{id}", todoController.Show).Methods("GET")

	// swagger:route POST /todos todos createTodo
	// Create a new todo
	// responses:
	// 	201: Todo
	// 422: ValidationError
	// 500: InternalServerError
	router.Post("/todos", todoController.Create).Name("CreateTodo")
	// apiRouter.HandleFunc("/todos", todoController.Create).Methods("POST")
	// swagger:route PUT /todos/{id} todos updateTodo
	// Update a todo
	// responses:
	// 200: Todo
	// 422: ValidationError
	// 500: InternalServerError
	router.Put("/todos/{id}", todoController.Update).Name("UpdateTodo")
	// apiRouter.HandleFunc("/todos/{id}", todoController.Update).Methods("PUT")

	// swagger:route DELETE /todos/{id} todos deleteTodo
	// Delete a todo
	// responses:
	// 200: Todo
	// 500: InternalServerError
	// 422: ValidationError
	// 404: NotFound
	// 401: Unauthorized
	router.Delete("/todos/{id}", todoController.Delete).Name("DeleteTodo")
	// apiRouter.HandleFunc("/todos/{id}", todoController.Delete).Methods("DELETE")
}

func registerAuthRoutes(apiRouter *router.Router) {
	authController := controllers.AuthController{}

	apiRouter.Post("/register", authController.Register).Name("Register")
	apiRouter.Post("/login", authController.Login).Name("Register")
	// apiRouter.HandleFunc("/register", authController.Register).Methods("POST")
	// apiRouter.HandleFunc("/login", authController.Login).Methods("POST")

}
