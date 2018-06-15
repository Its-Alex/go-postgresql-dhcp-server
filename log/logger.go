package log

import (
	"github.com/sirupsen/logrus"
)

// Logger is used to log
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	Logger.Formatter = &logrus.JSONFormatter{}
}

func ToggleVerbose(mode bool) {
	Logger.SetLevel(logrus.DebugLevel)
}
