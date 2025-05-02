package actions

import (
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/response"
	"imohamedsheta/gocrud/query"
	"net/http"
)

type TodoAction struct{}

func (t *TodoAction) CreateTodoAction(createTodo) (map[string]any, error) {
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

}
