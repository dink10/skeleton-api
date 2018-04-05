package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/skeleton-api/config"
)

// Server represents the core struct of this service
type Server struct {
	config *config.TemplateConfig
}

func main() {
	cfg := config.LoadConfig()
	setupLogger(cfg)

	srv := &Server{config: cfg}

	router := newRouter(srv)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}
