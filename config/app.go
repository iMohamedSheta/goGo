package config

import "imohamedsheta/gocrud/pkg/config"

func LoadAppConfig() {
	config.AppConfig.Set("app", map[string]any{
		"name": config.Env("APP_NAME", "GoCrudRestApi"),
		"port": config.Env("APP_PORT", 7777),
		"url":  config.Env("APP_URL", "http://localhost"),
	})
}
