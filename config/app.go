package config

import (
	"time"

	"github.com/iMohamedSheta/xapp/pkg/config"
)

func init() {
	config.Register(appConfig)
}

func appConfig() {
	config.App.Set("app", map[string]any{
		"name":             config.Env("APP_NAME", "GoCrudRestApi"),
		"url":              config.Env("APP_URL", "localhost"),
		"port":             config.Env("APP_PORT", 7777),
		"shutdown_timeout": 20 * time.Second,
		"log_path":         config.Env("APP_LOG_PATH", "storage/logs/app.json"),
		"env":              config.Env("APP_ENV", "dev"),
		"debug":            config.Env("APP_DEBUG", true),
		"secret":           config.Env("APP_SECRET", "hxdCTfhtkyJBVE01k8vvtaMHbzTmr401QqGl1111"),
		"jwt_expiry":       config.Env("APP_JWT_EXPIRY", 120),
		"auth": map[string]any{
			"type":                 "jwt",
			"access_token_expiry":  15 * time.Hour,
			"refresh_token_expiry": 168 * time.Hour,
		},
		"grpc": map[string]any{
			"port":    config.Env("GRPC_PORT", 50051),
			"host":    config.Env("GRPC_HOST", "localhost"),
			"enabled": config.Env("GRPC_ENABLED", true),
		},
	})
}
