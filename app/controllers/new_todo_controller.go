package controllers

import (
	"imohamedsheta/gocrud/app/enums"
	"imohamedsheta/gocrud/pkg/response"
	"imohamedsheta/gocrud/query"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoController struct{}

// return paginated list of todos of authenticated the user
func (c *TodoController) Index(w http.ResponseWriter, r *http.Request) {
	userIdRaw := r.Context().Value(enums.ContextKeyUserId)
	userId, ok := userIdRaw.(float64)

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

func (c *TodoController) Create(w http.ResponseWriter, r *http.Request) {

}

func (c *TodoController) Update(w http.ResponseWriter, r *http.Request) {

}

func (c *TodoController) Delete(w http.ResponseWriter, r *http.Request) {

}
