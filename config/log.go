package config

import "imohamedsheta/gocrud/pkg/config"

func LoadLogConfig() {
	config.App.Set("log", map[string]any{
		"default": "app",

		"channels": map[string]any{
			"app": map[string]any{
				"driver":   "daily",
				"path":     config.Env("APP_LOG_PATH", "storage/logs/app.log"),
				"level":    "debug",
				"max_size": 100, // in MB
				"max_age":  30,  // in days
				"backup":   false,
			},
			"dev": map[string]any{
				"driver":   "stdout",
				"level":    "debug",
				"max_size": 100, // in MB
				"max_age":  30,  // in days
				"backup":   false,
			},
		},
	})
}
