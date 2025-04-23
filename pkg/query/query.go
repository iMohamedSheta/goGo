package query

import (
	"errors"
	"imohamedsheta/gocrud/pkg/logger"
	"reflect"
	"slices"
)

func GetStructFieldNames(tableStruct any) []string {
	var fieldNames []string

	typ := reflect.TypeOf(tableStruct)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		column := typ.Field(i).Tag.Get("db")
		if column != "" {
			fieldNames = append(fieldNames, column)
		}
	}

	return fieldNames
}

func ValidateColumns(rowStruct any, columns []string) ([]string, error) {
	var selectedColumns []string
	allowed := GetStructFieldNames(rowStruct)

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
