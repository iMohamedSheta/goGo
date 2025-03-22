package config

import "os"

type Config []map[string]any

func (c *Config) Env(value string, defaultValue string) string {
	v := os.Getenv(value)
	if v == "" {
		return defaultValue
	}

	return v
}
