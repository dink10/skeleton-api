package logger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// InitLogger initialises logger
func Init(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse log level: %s, %v", logLevel, err))
	}
	log.SetLevel(lvl)
}
