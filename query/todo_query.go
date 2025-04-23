package query

import (
	"fmt"
	"imohamedsheta/gocrud/database"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/query"
	"strings"
	"time"
)

type Todo struct {
	Id          *int      `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      int8      `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type todosTable []Todo

func TodosTable() *todosTable {
	return &todosTable{}
}

func (t *todosTable) GetSql(columns ...string) (string, error) {

	if len(columns) == 0 || (len(columns) == 1 && columns[0] == "*") {
		return "SELECT * FROM todos", nil
	}

	validColumns, err := query.ValidateColumns(&Todo{}, columns)

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

	validColumns, err := query.ValidateColumns(&Todo{}, columns)

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

func (todos *todosTable) Insert(todo *Todo) error {
	db := database.DB()

	_, err := db.NamedExec(`
		INSERT INTO todos (title, description, status, created_at, updated_at)
		VALUES (:title, :description, :status, :created_at, :updated_at)
	`, todo)

	if err != nil {
		logger.Log().Error("Insert error: " + err.Error())
		return err
	}

	return nil
}
