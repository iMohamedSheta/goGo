package routes

import (
	"github.com/iMohamedSheta/xapp/app/controllers"
	"github.com/iMohamedSheta/xapp/app/middlewares"
	"github.com/iMohamedSheta/xapp/pkg/router"

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

	router.Post("/register", authController.Register).Name("register")
	router.Post("/login", authController.Login).Name("login")
	router.Post("/refresh/access-token", authController.RefreshAccessToken).Name("refresh_access_token")
}

func registerTodoRoutes(router *router.Router) {
	todoController := controllers.TodoController{}

	router.Get("/todos", todoController.Index).Name("list_todos")

	router.Get("/todos/{id}", todoController.Show).Name("show_todo")

	router.Post("/todos", todoController.Create).Name("create_todo")

	router.Put("/todos/{id}", todoController.Update).Name("update_todo")

	router.Delete("/todos/{id}", todoController.Delete).Name("delete_todo")
}
