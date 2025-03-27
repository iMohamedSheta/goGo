package support

import "imohamedsheta/gocrud/pkg/config"

// Config mimics Laravel's config() helper
func Config(path string) interface{} {
	return config.AppConfig.Get(path)
}
