package api

import (
	"net/http"
	"os"

	"bitbucket.org/gismart/{{Name}}/config"
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

	// middleware that should be applied to all requests
	return handlers.CombinedLoggingHandler(os.Stdout, handlers.RecoveryHandler()(router))
}
