package rules

import (
	"imohamedsheta/gocrud/database"
	"imohamedsheta/gocrud/pkg/logger"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Unique(fl validator.FieldLevel) bool {
	db := database.DB()
	param := fl.Param() // ex: "users,username"
	parts := strings.Split(param, ",")

	if len(parts) != 2 {
		return false
	}

	tableName := parts[0]
	columnName := parts[1]

	// It can be sql injected but i will trust the user cause the user is the developer
	// so he can add validation like this:
	// unique=users,email
	query := "SELECT COUNT(*) FROM " + tableName + " WHERE " + columnName + " = ?"

	var count int
	err := db.Get(&count, query, fl.Field().String())

	if err != nil {
		logger.Log().Error("error while checking unique constraint: " + err.Error())
		return false
	}

	return count == 0
}
