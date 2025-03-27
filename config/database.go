package config

import "imohamedsheta/gocrud/pkg/config"

func LoadDatabaseConfig() {
	config.AppConfig.Set("database", map[string]any{

		// This is the default database connection should be valid connection to use.
		"default": config.Env("DB_CONNECTION", "mysql"),

		"connections": map[string]any{

			// Mysql connection
			"mysql": map[string]any{
				"host":     config.Env("DB_HOST", "localhost"),
				"port":     config.Env("DB_PORT", "3306"),
				"user":     config.Env("DB_USERNAME", "root"),
				"pass":     config.Env("DB_PASSWORD", ""),
				"database": config.Env("DB_DATABASE", "go"),
				"charset":  "utf8",
				"driver":   "mysql",
			},

			// Postgres connection
			"pgsql": map[string]any{
				"host":     config.Env("DB_HOST", "localhost"),
				"port":     config.Env("DB_PORT", 5432),
				"user":     config.Env("DB_USERNAME", "root"),
				"pass":     config.Env("DB_PASSWORD", ""),
				"database": config.Env("DB_DATABASE", "go"),
				"charset":  "utf8",
				"driver":   "pgsql",
				"sslmode":  "disable",
				"timezone": "UTC+2",
				"prefix":   "",
				"singular": false,
				"schema":   "",
			},

			"sqlite": map[string]any{
				"driver":   "sqlite",
				"database": config.Env("DB_DATABASE", "db"+config.Env("APP_NAME").(string)),
				"prefix":   "",
				"singular": false,
			},
		},
	})
}
