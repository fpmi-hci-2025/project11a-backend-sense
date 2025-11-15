package logger

import (
	"github.com/sirupsen/logrus"
)

// New creates a new logger instance
func New() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	return log
}

