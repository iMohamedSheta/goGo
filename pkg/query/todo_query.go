package query

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xapp/app/models"
	"github.com/iMohamedSheta/xapp/database"
	"github.com/iMohamedSheta/xapp/pkg/logger"
)

type todosTable []models.Todo

func TodosTable() *todosTable {
	return &todosTable{}
}

func (t *todosTable) GetSql(columns ...string) (string, error) {

	if len(columns) == 0 || (len(columns) == 1 && columns[0] == "*") {
		return "SELECT * FROM todos", nil
	}

	validColumns, err := ValidateColumns(&models.Todo{}, columns)

	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf("SELECT %s FROM todos", strings.Join(validColumns, ", "))

	return sql, nil
}

func (todos *todosTable) Get(columns ...string) (*todosTable, error) {
	sql, err := todos.GetSql(columns...)

	if err != nil {
		logger.Log().Error(err.Error())
		return nil, err
	}

	db := database.DB()

	err = db.Select(todos, sql)

	if err != nil {
		logger.Log().Error(err.Error())
		return nil, err
	}

	return todos, nil
}

func (todos *todosTable) GetByIdSql(columns ...string) (string, error) {
	if len(columns) == 0 || (len(columns) == 1 && columns[0] == "*") {
		return "SELECT * FROM todos WHERE id = ?", nil
	}

	validColumns, err := ValidateColumns(&models.Todo{}, columns)

	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf("SELECT %s FROM todos WHERE id = ?", strings.Join(validColumns, ", "))

	return sql, nil
}

func (todos *todosTable) GetById(id int, columns ...string) (*todosTable, error) {
	sql, err := todos.GetByIdSql(columns...)

	if err != nil {
		logger.Log().Error(err.Error())
		return nil, err
	}

	db := database.DB()

	err = db.Get(todos, sql, id)

	if err != nil {
		logger.Log().Error(err.Error())
		return nil, err
	}

	return todos, nil
}

func (todos *todosTable) Insert(todo *models.Todo) error {
	db := database.DB()

	_, err := db.NamedExec(`
		INSERT INTO todos (title, user_id, description, status, created_at, updated_at)
		VALUES (:title, :user_id, :description, :status, :created_at, :updated_at)
	`, todo)

	if err != nil {
		logger.Log().Error("Insert error: " + err.Error())
		return err
	}

	return nil
}
