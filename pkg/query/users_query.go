package query

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xapp/app/models"
	"github.com/iMohamedSheta/xapp/database"
	"github.com/iMohamedSheta/xapp/pkg/logger"
)

type usersTable []models.User

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

	validColumns, err := ValidateColumns(models.User{}, columns)

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

func (users *usersTable) Insert(user *models.User) error {
	db := database.DB()

	_, err := db.NamedExec(`
		INSERT INTO users (username, email, password, first_name, last_name, created_at, updated_at)
		VALUES (:username, :email, :password, :first_name, :last_name, :created_at, :updated_at)
	`, user)

	if err != nil {
		logger.Log().Error("Insert error: " + err.Error())
		return err
	}

	return nil
}

func (users *usersTable) Update(user models.User) error {
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
