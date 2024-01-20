package logging

import (
	"career-compass-go/config"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitializeLogger creates a new instance for the global logger
func InitializeLogger() {
	Logger = logrus.New()

	if config.ViperConfig.GetString("RUN_MODE") == "release" {
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		Logger.SetLevel(logrus.DebugLevel)
	}

	Logger.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "2006-01-02T15:04:05.000"}
}
