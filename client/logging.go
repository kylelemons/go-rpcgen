package client

import (
	"log"
)

var (
	// Set this to something non-null if you want log messages to flow.
	Logger *log.Logger
)

func logMessage(format string, args ...interface{}) {
	if Logger == nil {
		return
	}

	Logger.Printf(format, args...)
}
