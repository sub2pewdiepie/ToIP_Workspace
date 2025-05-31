package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	// Set format based on env (default JSON)
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	if format == "text" {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch level {
	case "trace":
		Logger.SetLevel(logrus.TraceLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.DebugLevel)
	}
}
