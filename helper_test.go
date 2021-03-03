package main

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "badvalue")
	response := logLevel()
	if response != log.InfoLevel {
		t.Error()
	}
}
