package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type adapter func(http.Handler) http.Handler

var adapters []adapter

// all middleware should be set here
func init() {
	adapters = []adapter{loggerHandler(), handlers.CompressHandler}
}

func newRouter(srv *Server) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	healthCheckHandlers(router)

	return router
}

// wrapMiddleware wraps handler with all specified middlewares
func wrapMiddleware(f http.HandlerFunc) http.Handler {
	return adapt(f)
}

func adapt(h http.Handler) http.Handler {
	for _, currentAdapter := range adapters {
		h = currentAdapter(h)
	}
	return h
}

func healthCheckHandlers(router *mux.Router) {
	router.Methods("GET").Path("/index").Handler(wrapMiddleware(healthCheck))
}
