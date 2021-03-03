package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func logLevel() log.Level {
	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch value {
		case "trace":
			return log.TraceLevel
		case "debug":
			return log.DebugLevel
		case "info":
			return log.InfoLevel
		case "warn":
			return log.WarnLevel
		case "error":
			return log.ErrorLevel
		case "fatal":
			return log.FatalLevel
		case "panic":
			return log.WarnLevel
		default:
			return log.InfoLevel
		}
	}
	return log.InfoLevel
}
