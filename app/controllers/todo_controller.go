package controllers

import (
	"encoding/json"
	"imohamedsheta/gocrud/app/enums"
	"imohamedsheta/gocrud/app/requests"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/response"
	"imohamedsheta/gocrud/pkg/validate"
	"imohamedsheta/gocrud/query"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type TodoController struct{}

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

	userTodos, meta, err := query.Table("todos").Where("user_id", "=", int64(userId)).Paginate(page, perPage, true)

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

	itemId, err := strconv.Atoi(itemIdRaw)

	if err != nil {
		response.ErrorJson(w, "Invalid id parameter", "invalid_id", http.StatusBadRequest)
		return
	}

	todo, err := query.Table("todos").Where("id", "=", itemId).First()

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

	// Validate request
	ok, validationErrors := validate.ValidateRequest(&req)

	if !ok {
		response.ValidationErrorJson(w, validationErrors)
		return
	}

	createdTodo := []map[string]any{
		{
			"title":       req.Title,
			"description": req.Description,
			"user_id":     userId,
			"status":      uint8(enums.CANCELLED),
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
	}

	sqlResult, err := query.Table("todos").Insert(createdTodo)

	if err != nil {
		logger.Log().Error(err.Error())
		response.ErrorJson(w, "Failed to create todo", "failed_to_create_todo", http.StatusInternalServerError)
		return
	}

	createdTodoId, err := sqlResult.LastInsertId()

	if err == nil {
		createdTodo[0]["id"] = createdTodoId
	}

	data := map[string]any{
		"todo": createdTodo[0],
	}

	response.Json(w, "Todo created successfully", data, http.StatusCreated)
}

func (c *TodoController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemRawId := vars["id"]

	itemId, err := strconv.Atoi(itemRawId)

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

	sqlResult, err := query.Table("todos").Where("id", "=", itemId).Where("user_id", "=", userId).Update(map[string]any{
		"title":       req.Title,
		"description": req.Description,
		"updated_at":  time.Now(),
	})

	if err != nil {
		response.ErrorJson(w, "Error updating todo", "error_updating_todo", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := sqlResult.RowsAffected()

	if err != nil {
		response.ServerErrorJson(w)
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

func (c *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemIdRaw := vars["id"]

	itemId, err := strconv.Atoi(itemIdRaw)

	if err != nil {
		response.ErrorJson(w, "Invalid item id", "invalid_item_id", http.StatusBadRequest)
		return
	}

	userId, ok := getUserIdFromContext(r)

	if !ok {
		response.ErrorJson(w, "Unauthorized Action", "unauthenticated", http.StatusUnauthorized)
		return
	}

	sqlResult, err := query.Table("todos").Where("user_id", "=", userId).Where("id", "=", itemId).Delete()

	if err != nil {
		response.ErrorJson(w, "Error deleting todo", "error_deleting_todo", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := sqlResult.RowsAffected()

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
