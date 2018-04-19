package api

import (
	"bitbucket.org/gismart/skeleton-api/config"
	"github.com/gorilla/mux"
)

// Server represents the core struct of this service
type Server struct {
	config *config.APIConfig
}

// InitializeAPI initializes API service with all necessary routes and returns configured router
func InitializeAPI(config *config.APIConfig) *mux.Router {
	srv := &Server{config: config}
	router := mux.NewRouter().StrictSlash(true)

	addBasicRoutes(srv, router)

	return router
}
