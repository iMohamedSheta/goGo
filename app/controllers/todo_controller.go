package controllers

import (
	"encoding/json"
	"fmt"
	"imohamedsheta/gocrud/pkg/config"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/query"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoController struct {
}

func (c *TodoController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := query.UsersTable().Get()

	if err != nil {
		logger.Log().Error(err.Error())
		errorResponse(w, "Error getting users")
		return
	}

	message := "Welcome to the " + config.App.Get("app.name").(string) + " API"
	response := map[string]any{
		"users":   users,
		"message": message,
	}

	responseJson, err := json.Marshal(response)

	if err != nil {
		logger.Log().Error(err.Error())
		errorResponse(w, "Error marshalling users")
		return
	}

	w.Write(responseJson)
}

func (c *TodoController) Show(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid ID"))
		return
	}

	w.Write([]byte(fmt.Sprintf("Show todo with ID: %d", id)))
}

func (c *TodoController) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte("Create todo with id " + fmt.Sprintf("%v", vars)))
}

func (c *TodoController) Update(w http.ResponseWriter, r *http.Request) {

}

func (c *TodoController) Delete(w http.ResponseWriter, r *http.Request) {

}

func errorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
