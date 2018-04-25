package api

import (
	"net/http"

	"bitbucket.org/gismart/skeleton-api/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server represents the core struct of this service
type Server struct {
	config *config.APIConfig
}

// InitializeAPI initializes API service with all necessary routes and returns configured router
func InitializeAPI(config *config.APIConfig) http.Handler {
	srv := &Server{config: config}
	router := mux.NewRouter().StrictSlash(true)

	addBasicRoutes(srv, router)

	return handlers.RecoveryHandler()(router)
}
