package actions

import (
	"net/http"

	"github.com/iMohamedSheta/xapp/pkg/logger"
	"github.com/iMohamedSheta/xapp/pkg/response"
	"github.com/iMohamedSheta/xapp/query"
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
