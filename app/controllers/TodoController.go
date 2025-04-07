package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoController struct {
}

func (c *TodoController) Index(w http.ResponseWriter, r *http.Request) {
	message := "Welcome to the " + os.Getenv("APP_NAME") + " API"

	w.Write([]byte(message))
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
