package utils

import (
	"github.com/sirupsen/logrus"
)

// GetLogLevel log level
func GetLogLevel(level string) (logLevel logrus.Level) {
	defer func() {
		logrus.Infof("Set log level to %s", logLevel)
	}()
	if len(level) == 0 {
		logLevel = logrus.InfoLevel
		return
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	return
}
