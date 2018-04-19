package main

import (
	"os"

	"bitbucket.org/gismart/skeleton-api/config"
	log "github.com/sirupsen/logrus"
)

func setupLogger(config *config.APIConfig) {
	log.SetOutput(os.Stdout)
	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("Failed to parse log level. %v", err)
	}
	log.SetLevel(lvl)
}
