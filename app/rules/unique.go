package rules

import (
	"strings"

	"github.com/iMohamedSheta/xapp/database"
	"github.com/iMohamedSheta/xapp/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func Unique(fl validator.FieldLevel) bool {
	db := database.DB()
	param := fl.Param() // ex: "users|username"
	parts := strings.Split(param, "_")

	if len(parts) != 2 {
		logger.Log().Error("Invalid unique validator format, expected: table|column")
		return false
	}

	tableName := strings.TrimSpace(parts[0])
	columnName := strings.TrimSpace(parts[1])

	if tableName == "" || columnName == "" {
		return false
	}

	query := "SELECT COUNT(*) FROM " + tableName + " WHERE " + columnName + " = ?"

	var count int
	err := db.Get(&count, query, fl.Field().String())

	if err != nil {
		logger.Log().Error("error while checking unique constraint: " + err.Error())
		return false
	}

	return count == 0
}
