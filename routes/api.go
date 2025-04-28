package routes

import (
	"imohamedsheta/gocrud/app/controllers"
	"imohamedsheta/gocrud/app/middlewares"
	"imohamedsheta/gocrud/pkg/router"

	"github.com/gorilla/mux"
)

var authMiddleware = []mux.MiddlewareFunc{
	middlewares.AuthMiddleware(),
	middlewares.JSONContentTypeMiddleware(),
}

func RegisterApiRoutes(r *router.Router) {

	registerAuthRoutes(r)

	r.Group("auth", "/", authMiddleware,
		func(router *router.Router) {
			registerTodoRoutes(router)
		})
}

func registerAuthRoutes(router *router.Router) {
	authController := controllers.AuthController{}

	router.Post("/register", authController.Register).Name("Register")
	router.Post("/login", authController.Login).Name("Login")
	router.Post("/refresh/access-token", authController.RefreshAccessToken).Name("RefreshAccessToken")
}

func registerTodoRoutes(router *router.Router) {
	todoController := controllers.TodoController{}

	router.Get("/todos", todoController.Index).Name("ListTodos")

	router.Get("/todos/{id}", todoController.Show).Name("ListTodos")

	router.Post("/todos", todoController.Create).Name("CreateTodo")

	router.Put("/todos/{id}", todoController.Update).Name("UpdateTodo")

	router.Delete("/todos/{id}", todoController.Delete).Name("DeleteTodo")
}
