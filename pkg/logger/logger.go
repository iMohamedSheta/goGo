package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/iMohamedSheta/xapp/pkg/enums"
	s "github.com/iMohamedSheta/xapp/pkg/support"

	"go.uber.org/zap"
)

var (
	// logger is the global logger instance
	App  *zap.Logger
	once sync.Once
)

// Load the logger
func LoadLogger() *zap.Logger {
	once.Do(func() {
		// Get the log path from the config file
		logPathRaw, _ := s.Config("app.log_path")

		logPath := logPathRaw.(string)

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
	})

	return App
}

// Get the logger instance
func Log() *zap.Logger {
	return App
}
