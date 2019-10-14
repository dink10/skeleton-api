package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Run(address string, withTracing bool) {
	router := runRoute(withTracing)

	log.Fatal(http.ListenAndServe(address, router))
}
