package config

import (
	"fmt"
	"imohamedsheta/gocrud/helpers"
	"os"
)

var Databases map[string]any

func init() {
	Databases = map[string]any{
		"connections": map[string]any{

			// This is the default database connection should be valid connection to use.
			"default": os.Getenv("DB_CONNECTION"),

			// Mysql connection
			"mysql": map[string]any{
				"host":    Config.Env("DB_HOST", "localhost"),
				"port":    os.Getenv("DB_PORT"),
				"user":    os.Getenv("DB_USER"),
				"pass":    os.Getenv("DB_PASS"),
				"db_name": os.Getenv("DB_NAME"),
				"charset": "utf8",
				"driver":  "mysql",
			},

			// Postgres connection
			"pgsql": map[string]any{
				"host":    os.Getenv("DB_HOST"),
				"port":    os.Getenv("DB_PORT"),
				"user":    os.Getenv("DB_USER"),
				"pass":    os.Getenv("DB_PASS"),
				"db_name": os.Getenv("DB_NAME"),
				"charset": "utf8",
				"driver":  "pgsql",
			},
		},
	}
}

// GetDefaultDatabaseConfig retrieves the active database configuration
func GetDefaultDatabaseConfig() map[string]any {
	connections := Databases["connections"].(map[string]any)
	helpers.LogError(connections)
	defaultDB, ok := connections["default"].(string)
	helpers.LogError("the default database connection is " + defaultDB)

	if !ok || defaultDB == "" {
		panic("Default database connection is not set or invalid")
	}

	config, exists := connections[defaultDB]
	if !exists {
		panic(fmt.Sprintf("Database configuration for '%s' not found", defaultDB))
	}

	return config.(map[string]any)
}
