package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()

	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Logger.WithField("error", err).Error("Failed to create logs directory")
	}

	// Open log file
	logFile, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.WithField("error", err).Error("Failed to open log file")
	} else {
		// Write to both console and file
		Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	// Set format based on env (default JSON)
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	if format == "text" {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// Set level based on env (default Debug for dev)
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
