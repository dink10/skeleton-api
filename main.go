package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/gismart/skeleton-api/api"
	"bitbucket.org/gismart/skeleton-api/config"
)

func main() {
	cfg := config.LoadConfig()
	setupLogger(cfg)

	router := api.InitializeAPI(cfg)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}
