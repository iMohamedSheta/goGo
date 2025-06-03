package support

import "github.com/iMohamedSheta/xapp/pkg/config"

// Config returns the configuration value for the given path example: support.Config("app.log_path")
func Config(path string) (any, error) {
	return config.App.Get(path)
}
