package config

import (
	"os"
	"strings"
)

// var (
// 	// configStore holds all configuration arrays
// 	configStore = make(map[string]any)
// )

type Config struct {
	store map[string]any
}

// Global app config instance
var App = &Config{store: make(map[string]any)}

// Set a configuration value
func (c *Config) Set(key string, value any) {
	c.store[key] = value
}

// Get a configuration value
func (c *Config) Get(key string) any {
	keys := strings.Split(key, ".")
	current := c.store
	lastIndex := (len(keys) - 1)

	for i, k := range keys {

		// if we are at the last key return the value
		if i == lastIndex {
			if value, ok := current[k]; ok {
				return value
			}
			panic("config key not found")
		}

		// check if the key inside nested map
		if value, ok := current[k].(map[string]any); ok {
			current = value
		} else {
			panic("config key not found")
		}
	}

	panic("config key not found")
}

// Load configuration from environment variables
func Env(key string, defaultValue ...any) any {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}
