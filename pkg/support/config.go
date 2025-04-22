package support

import "imohamedsheta/gocrud/pkg/config"

// Config returns the configuration value for the given path example: support.Config("app.log_path")
func Config(path string) interface{} {
	return config.App.Get(path)
}
