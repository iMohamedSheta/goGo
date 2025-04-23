package query

import (
	"errors"
	"fmt"
	"imohamedsheta/gocrud/database"
	"imohamedsheta/gocrud/pkg/logger"
	"reflect"
	"slices"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type usersTable []User

func UsersTable() *usersTable {
	return &usersTable{}
}

/*
* @return string - SQL query for selecting users from the database.
 */
func (u *usersTable) GetSql(columns ...string) string {

	if len(columns) == 0 || (len(columns) == 1 && columns[0] == "*") {
		return "SELECT * FROM users"
	}

	validColumns, err := validateColumns(columns)

	if err != nil {
		return ""
	}

	sql := fmt.Sprintf("SELECT %s FROM users", strings.Join(validColumns, ", "))

	return sql
}

func (users *usersTable) Get(columns ...string) (*usersTable, error) {
	sql := users.GetSql(columns...)
	db := database.DB()

	err := db.Select(users, sql)

	if err != nil {
		logger.Log().Error("Query error: " + err.Error())
		return nil, err
	}

	return users, nil
}

func (users *usersTable) Insert(user User) error {
	db := database.DB()

	_, err := db.NamedExec(`
		INSERT INTO users (name, email, created_at)
		VALUES (:name, :email, :created_at)
	`, user)

	if err != nil {
		logger.Log().Error("Insert error: " + err.Error())
		return err
	}

	return nil
}

func (users *usersTable) Update(user User) error {
	db := database.DB()

	_, err := db.NamedExec(`
		UPDATE users SET name = :name, email = :email
		WHERE id = :id
	`, user)

	if err != nil {
		logger.Log().Error("Update error: " + err.Error())
		return err
	}

	return nil
}

func getStructFieldNames(tableStruct any) []string {
	var fieldNames []string

	typ := reflect.TypeOf(tableStruct)

	for i := 0; i < typ.NumField(); i++ {
		column := typ.Field(i).Tag.Get("db")
		if column != "" {
			fieldNames = append(fieldNames, column)
		}
	}

	return fieldNames
}

func validateColumns(columns []string) ([]string, error) {
	var selectedColumns []string
	allowed := getStructFieldNames(User{})

	for _, column := range columns {
		if slices.Contains(allowed, column) {
			selectedColumns = append(selectedColumns, column)
		} else {
			logger.Log().Error("Invalid column name: " + column)
			return nil, errors.New("Invalid column name" + column)
		}
	}

	return selectedColumns, nil
}
