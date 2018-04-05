package main

import (
	"net/http"
	"os"
	"time"

	"bitbucket.org/gismart/skeleton-api/config"
	log "github.com/sirupsen/logrus"
)

// loggerHandler logs all HTTP requests to output in a proper format
func loggerHandler() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			log.Debugf("%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				time.Since(start),
			)
		})
	}
}

func setupLogger(config *config.TemplateConfig) {
	log.SetOutput(os.Stdout)
	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("Failed to parse log level. %v", err)
	}
	log.SetLevel(lvl)
}
