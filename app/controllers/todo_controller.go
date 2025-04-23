package controllers

import (
	"encoding/json"
	"fmt"
	"imohamedsheta/gocrud/app/enums"
	"imohamedsheta/gocrud/pkg/config"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/validate"
	"imohamedsheta/gocrud/query"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type TodoController struct {
}

func (c *TodoController) UsersIndex(w http.ResponseWriter, r *http.Request) {
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

func (c *TodoController) Index(w http.ResponseWriter, r *http.Request) {
	todo := &query.Todo{
		Title:       "todo number 2",
		Description: "this is the todo number 2",
		Status:      int8(enums.IN_PROGRESS),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := query.TodosTable().Insert(todo)

	if err != nil {
		logger.Log().Error(err.Error())
		errorResponse(w, err.Error())

		return
	}

	w.Write([]byte("Todo created"))
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
	var todo query.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		logger.Log().Error(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"title":       todo.Title,
		"description": todo.Description,
	}

	// Define the validation rules
	rules := map[string]string{
		"title":       "required,min=3,max=100",
		"description": "required,min=5",
	}

	// Define custom error messages
	messages := map[string]string{
		"title":       "The Title is required and must be between 3 and 100 characters.",
		"description": "The Description is required and must be at least 5 characters long.",
	}

	// Perform validation with custom messages
	ok, validationErrors := validate.Validate(data, rules, messages)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
		return
	}

	// Set default values for status and timestamps
	todo.Status = int8(enums.CANCELLED)
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	// Insert the todo into the database
	if err := query.TodosTable().Insert(&todo); err != nil {
		logger.Log().Error(err.Error())
		errorResponse(w, "Failed to create todo")
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Todo created successfully",
	})
}

func (c *TodoController) Update(w http.ResponseWriter, r *http.Request) {

}

func (c *TodoController) Delete(w http.ResponseWriter, r *http.Request) {

}

func errorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
