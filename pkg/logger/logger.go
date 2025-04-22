package logger

import (
	"imohamedsheta/gocrud/pkg/enums"
	s "imohamedsheta/gocrud/pkg/support"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// logger is the global logger instance
var App *zap.Logger

// Load the logger
func LoadLogger() *zap.Logger {
	// Get the log path from the config file
	logPath := s.Config("app.log_path").(string)

	// Ensure the directory exists
	logDir := filepath.Dir(logPath)
	if _, err := os.Stat(logDir); err != nil {
		if err := os.Mkdir(logDir, 0755); err != nil {
			log.Fatal(enums.Red.Value() + "Failed to create log directory: " + err.Error() + enums.Reset.Value())
		}
	}

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logPath}
	cfg.ErrorOutputPaths = []string{logPath}
	logger, err := cfg.Build()

	if err != nil {
		panic("Failed to build logger: " + err.Error())
	}

	App = logger

	return App
}

// Get the logger instance
func Log() *zap.Logger {
	return App
}
