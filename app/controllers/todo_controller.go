package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/iMohamedSheta/xapp/app/enums"
	"github.com/iMohamedSheta/xapp/app/repository"
	"github.com/iMohamedSheta/xapp/app/requests"
	"github.com/iMohamedSheta/xapp/pkg/logger"
	"github.com/iMohamedSheta/xapp/pkg/response"
	"github.com/iMohamedSheta/xapp/pkg/validate"

	"github.com/gorilla/mux"
)

type TodoController struct {
	repository.TodoRepository
}

// return paginated list of todos of authenticated the user
func (c *TodoController) Index(w http.ResponseWriter, r *http.Request) {
	userId, ok := getUserIdFromContext(r)

	if !ok {
		response.ErrorJson(w, "Unauthorized Action", "unauthenticated", http.StatusUnauthorized)
		return
	}

	perPageRaw := r.URL.Query().Get("per_page")
	pageRaw := r.URL.Query().Get("page")

	var perPage int
	var page int

	if perPageRaw == "" {
		perPage = 100
	} else {
		var err error
		perPage, err = strconv.Atoi(perPageRaw)
		if err != nil {
			response.ErrorJson(w, "Invalid per_page parameter", "invalid_per_page", http.StatusBadRequest)
			return
		}
	}

	if pageRaw == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(pageRaw)
		if err != nil {
			response.ErrorJson(w, "Invalid page parameter", "invalid_page", http.StatusBadRequest)
			return
		}
	}

	userTodos, meta, err := c.TodoRepository.PaginatedUserTodos(int64(userId), perPage, page, true)

	if err != nil {
		response.ErrorJson(w, "Error fetching todos", "error_fetching_todos", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"todos": userTodos,
		"meta":  meta,
	}

	response.Json(w, "success", data, http.StatusOK)
}

// Get a specific todo of authenticated the user
func (c *TodoController) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemIdRaw := vars["id"]

	userId, ok := getUserIdFromContext(r)

	if !ok {
		response.ErrorJson(w, "Unauthorized Action", "unauthenticated", http.StatusUnauthorized)
		return
	}

	itemId, err := strconv.Atoi(itemIdRaw)

	if err != nil {
		response.ErrorJson(w, "Invalid id parameter", "invalid_id", http.StatusBadRequest)
		return
	}

	todo, err := c.TodoRepository.Find(int64(userId), int64(itemId))

	// todo, err := query.Table("todos").Where("id", "=", itemId).First()

	if err != nil {
		response.ErrorJson(w, "Error fetching todo", "error_fetching_todo", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"todo": todo,
	}

	response.Json(w, "success", data, http.StatusOK)
}

// Create a new todo for the authenticated the user
func (c *TodoController) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log().Error(err.Error())
		response.ErrorJson(w, "Invalid request", "invalid_request", http.StatusBadRequest)
		return
	}

	userId, ok := getUserIdFromContext(r)

	if !ok {
		response.ErrorJson(w, "Unauthorized Action", "unauthenticated", http.StatusUnauthorized)
		return
	}

	ok, validationErrors := validate.ValidateRequest(&req)

	if !ok {
		response.ValidationErrorJson(w, validationErrors)
		return
	}

	createTodo := []map[string]any{
		{
			"title":       req.Title,
			"description": req.Description,
			"user_id":     userId,
			"status":      uint8(enums.IN_PROGRESS),
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
	}

	createdTodo, err := c.TodoRepository.Create(int64(userId), createTodo)

	if err != nil {
		logger.Log().Error(err.Error())
		response.ErrorJson(w, "Failed to create todo", "failed_to_create_todo", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"todo": createdTodo,
	}

	response.Json(w, "Todo created successfully", data, http.StatusCreated)
}

// Update a todo for the authenticated user
func (c *TodoController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemIdRaw := vars["id"]

	itemId, err := strconv.ParseInt(itemIdRaw, 10, 64)

	if err != nil {
		response.ErrorJson(w, "Invalid id", "invalid_id", http.StatusBadRequest)
		return
	}

	userId, ok := getUserIdFromContext(r)

	if !ok {
		response.ErrorJson(w, "Unauthorized Action", "unauthenticated", http.StatusUnauthorized)
		return
	}

	req := requests.UpdateTodoRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJson(w, "Invalid request", "invalid_request", http.StatusBadRequest)
		return
	}

	ok, validationErrors := validate.ValidateRequest(&req)

	if !ok {
		response.ValidationErrorJson(w, validationErrors)
		return
	}

	updatedFields := map[string]any{
		"title":       req.Title,
		"description": req.Description,
	}

	rowsAffected, err := c.TodoRepository.Update(int64(userId), itemId, updatedFields)

	if err != nil {
		response.ErrorJson(w, "Error updating todo", "error_updating_todo", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		response.ErrorJson(w, "Todo not found", "todo_not_found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"todo": map[string]any{
			"id":          itemId,
			"title":       req.Title,
			"description": req.Description,
			"updated_at":  time.Now(),
		},
	}

	response.Json(w, "Todo updated successfully", data, http.StatusOK)
}

// Delete a specific todo of authenticated the user
func (c *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemIdRaw := vars["id"]

	itemId, err := strconv.ParseInt(itemIdRaw, 10, 64)

	if err != nil {
		response.ErrorJson(w, "Invalid item id", "invalid_item_id", http.StatusBadRequest)
		return
	}

	userId, ok := getUserIdFromContext(r)

	if !ok {
		response.ErrorJson(w, "Unauthorized Action", "unauthenticated", http.StatusUnauthorized)
		return
	}

	rowsAffected, err := c.TodoRepository.Delete(int64(userId), itemId)

	if err != nil {
		response.ErrorJson(w, "Error deleting todo", "error_deleting_todo", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		response.ErrorJson(w, "Todo not found", "todo_not_found", http.StatusNotFound)
		return
	}

	response.Json(w, "Todo deleted successfully", nil, http.StatusOK)
}

func getUserIdFromContext(r *http.Request) (float64, bool) {
	userIdRaw := r.Context().Value(enums.ContextKeyUserId)
	userId, ok := userIdRaw.(float64)

	return userId, ok
}
