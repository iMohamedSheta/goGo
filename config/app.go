package config

import "imohamedsheta/gocrud/pkg/config"

func LoadAppConfig() {
	config.App.Set("app", map[string]any{
		"name":     config.Env("APP_NAME", "GoCrudRestApi"),
		"port":     config.Env("APP_PORT", 7777),
		"url":      config.Env("APP_URL", "http://localhost"),
		"log_path": config.Env("APP_LOG_PATH", "storage/logs/app.log"),
		"env":      config.Env("APP_ENV", "dev"),
		"debug":    config.Env("APP_DEBUG", true),
	})
}
